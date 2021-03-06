package appmesh

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"unicode/utf8"

	"github.com/solo-io/solo-kit/pkg/utils/nameutils"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	"github.com/solo-io/supergloo/pkg/translator/utils"

	"github.com/mitchellh/hashstructure"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/appmesh"

	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/supergloo/pkg/api/v1"
)

// NOTE: copy-pasted from discovery/pkg/fds/discoveries/aws/aws.go
// TODO: aggregate these somewhere
const (
	// expected map identifiers for secrets
	awsAccessKey = "access_key"
	awsSecretKey = "secret_key"
)

func MeshName(meshRef core.ResourceRef) string {
	return fmt.Sprintf("%v-%v", meshRef.Namespace, meshRef.Name)
}

// TODO: util method
func NewAwsClientFromSecret(awsSecret *gloov1.Secret_Aws, region string) (*appmesh.AppMesh, error) {
	accessKey := awsSecret.Aws.AccessKey
	if accessKey != "" && !utf8.Valid([]byte(accessKey)) {
		return nil, errors.Errorf("%s not a valid string", awsAccessKey)
	}
	secretKey := awsSecret.Aws.SecretKey
	if secretKey != "" && !utf8.Valid([]byte(secretKey)) {
		return nil, errors.Errorf("%s not a valid string", awsSecretKey)
	}
	sess, err := session.NewSession(aws.NewConfig().
		WithCredentials(credentials.NewStaticCredentials(accessKey, secretKey, "")))
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create AWS session")
	}
	svc := appmesh.New(sess, &aws.Config{Region: aws.String(region)})

	return svc, nil
}

// todo: replace with interface
type AppMeshClient = *appmesh.AppMesh

type AppMeshSyncer struct {
	lock           sync.Mutex
	activeSessions map[uint64]AppMeshClient
}

func NewSyncer() *AppMeshSyncer {
	return &AppMeshSyncer{
		lock:           sync.Mutex{},
		activeSessions: make(map[uint64]AppMeshClient),
	}
}

func hashCredentials(awsSecret *gloov1.Secret_Aws, region string) uint64 {
	hash, _ := hashstructure.Hash(struct {
		awsSecret *gloov1.Secret_Aws
		region    string
	}{
		awsSecret: awsSecret,
		region:    region,
	}, nil)
	return hash
}

func (s *AppMeshSyncer) Sync(ctx context.Context, snap *v1.TranslatorSnapshot) error {
	ctx = contextutils.WithLogger(ctx, "appmesh-syncer")
	logger := contextutils.LoggerFrom(ctx)
	meshes := snap.Meshes.List()
	upstreams := snap.Upstreams.List()
	secrets := snap.Secrets.List()
	rules := snap.Routingrules.List()

	logger.Infof("begin sync %v (%v meshes, %v upstreams, %v rules, %v secrets)", snap.Hash(),
		len(meshes), len(upstreams), len(rules), len(secrets))
	defer logger.Infof("end sync %v", snap.Hash())
	logger.Debugf("%v", snap)

	for _, mesh := range snap.Meshes.List() {
		if err := s.sync(ctx, mesh, snap); err != nil {
			return errors.Wrapf(err, "syncing mesh %v", mesh.Metadata.Ref())
		}
	}

	return nil
	/*
		0 - mesh per mesh
		1 - virtual node per upstream
		2 - routing rules get aggregated into virtual service like object
		routes on virtual service become the aws routes
		virtual service becomes virtual router
	*/
	//exampleMesh := appmesh.CreateMeshInput{}
	//exampleVirtualNode := appmesh.CreateVirtualNodeInput{}
	//exampleVirtualRouter := appmesh.CreateVirtualRouterInput{}
	//exampleRoute := appmesh.CreateRouteInput{}
	//
	//destinationRules, err := virtualNodesForUpstreams(rules, meshes, upstreams)
	//if err != nil {
	//	return errors.Wrapf(err, "creating subsets from snapshot")
	//}
	//
	//virtualServices, err := virtualServicesForRules(rules, meshes, upstreams)
	//if err != nil {
	//	return errors.Wrapf(err, "creating virtual services from snapshot")
	//}
	//return s.writeappmeshCrds(ctx, destinationRules, virtualServices)
}

func (s *AppMeshSyncer) sync(ctx context.Context, mesh *v1.Mesh, snap *v1.TranslatorSnapshot) error {
	appMesh, ok := mesh.MeshType.(*v1.Mesh_AppMesh)
	if !ok {
		return nil
	}
	if appMesh.AppMesh == nil {
		return errors.Errorf("%v missing configuration for AppMesh", mesh.Metadata.Ref())
	}
	meshName := MeshName(mesh.Metadata.Ref())
	desiredMesh := appmesh.CreateMeshInput{
		MeshName: aws.String(meshName),
	}
	upstreams := snap.Upstreams.List()
	routingRules := snap.Routingrules.List()

	virtualNodes, virtualRouters, routes, err := desiredResources(meshName, upstreams, routingRules)
	if err != nil {
		return errors.Wrapf(err, "generating desired resources")
	}

	secrets := snap.Secrets.List()
	client, err := s.newOrCachedClient(appMesh.AppMesh, secrets)
	if err != nil {
		return errors.Wrapf(err, "creating new AWS AppMesh session")
	}

	logger := contextutils.LoggerFrom(ctx)
	logger.Infof("syncing desired state")
	logger.Debugf("desired mesh: %v", desiredMesh)
	logger.Debugf("desired virtual nodes: %v", virtualNodes)
	logger.Debugf("desired virtual routers: %v", virtualRouters)
	logger.Debugf("desired routes: %v", routes)
	if err := resyncState(client, desiredMesh, virtualNodes, virtualRouters, routes); err != nil {
		return errors.Wrapf(err, "reconciling desired state")
	}
	return nil
}

func (s *AppMeshSyncer) newOrCachedClient(appMesh *v1.AppMesh, secrets gloov1.SecretList) (AppMeshClient, error) {
	secret, err := secrets.Find(appMesh.AwsCredentials.Strings())
	if err != nil {
		return nil, errors.Wrapf(err, "finding aws credentials for mesh")
	}
	region := appMesh.AwsRegion
	if region == "" {
		return nil, errors.Wrapf(err, "mesh must provide aws_region")
	}

	awsSecret, ok := secret.Kind.(*gloov1.Secret_Aws)
	if !ok {
		return nil, errors.Errorf("mesh referenced non-AWS secret, AWS secret required")
	}
	if awsSecret.Aws == nil {
		return nil, errors.Errorf("secret missing field Aws")
	}

	// check if we already have an active session for this region/credential
	sessionKey := hashCredentials(awsSecret, region)
	s.lock.Lock()
	appMeshClient, ok := s.activeSessions[sessionKey]
	s.lock.Unlock()
	if !ok {
		// create a new client and cache it
		// TODO: is there a point where we should drop old sessions?
		// maybe aws will do it for us
		appMeshClient, err = NewAwsClientFromSecret(awsSecret, region)
		if err != nil {
			return nil, errors.Wrapf(err, "creating aws client from provided secret/region")
		}
		s.lock.Lock()
		s.activeSessions[sessionKey] = appMeshClient
		s.lock.Unlock()
	}

	return appMeshClient, nil
}

func desiredResources(meshName string, upstreams gloov1.UpstreamList, routingRules v1.RoutingRuleList) ([]appmesh.CreateVirtualNodeInput, []appmesh.CreateVirtualRouterInput, []appmesh.CreateRouteInput, error) {
	virtualNodes, err := virtualNodesFromUpstreams(meshName, upstreams)
	if err != nil {
		return nil, nil, nil, errors.Wrapf(err, "creating virtual nodes from upstreams")
	}
	virtualRouters, err := virtualRoutersFromUpstreams(meshName, upstreams)
	if err != nil {
		return nil, nil, nil, errors.Wrapf(err, "creating virtual routers from upstreams")
	}
	routes, err := routesForUpstreams(meshName, upstreams, routingRules)
	if err != nil {
		return nil, nil, nil, errors.Wrapf(err, "creating virtual routers from upstreams")
	}
	return virtualNodes, virtualRouters, routes, nil
}

func resyncState(client AppMeshClient,
	mesh appmesh.CreateMeshInput,
	vNodes []appmesh.CreateVirtualNodeInput,
	vRouters []appmesh.CreateVirtualRouterInput,
	routes []appmesh.CreateRouteInput) error {
	if err := reconcileMesh(client, mesh); err != nil {
		return errors.Wrapf(err, "reconciling mesh")
	}
	if err := reconcileVirtualNodes(client, mesh, vNodes); err != nil {
		return errors.Wrapf(err, "reconciling virtual nodes")
	}
	if err := reconcileVirtualRouters(client, mesh, vRouters); err != nil {
		return errors.Wrapf(err, "reconciling virtual routers")
	}
	for _, vRouter := range vRouters {
		if err := reconcileRoutes(client, mesh, *vRouter.VirtualRouterName, routes); err != nil {
			return errors.Wrapf(err, "reconciling routes")
		}
	}

	return nil
}

func getUniqueHosts(upstreams gloov1.UpstreamList) ([]string, error) {
	var uniqueHosts []string
	for _, us := range upstreams {
		host, err := getHostForUpstream(us)
		if err != nil {
			return nil, errors.Wrapf(err, "getting host for upstream")
		}
		// only add unique ports
		var alreadyAdded bool
		for _, addedHost := range uniqueHosts {
			if addedHost == host {
				alreadyAdded = true
				break
			}
		}
		if !alreadyAdded {
			uniqueHosts = append(uniqueHosts, host)
		}
	}
	return uniqueHosts, nil
}

// Unused currently...
func getAllHosts(upstreams gloov1.UpstreamList) ([]string, error) {
	uniqueHostMap := make(map[string]struct{})
	var uniqueHosts []string
	for _, us := range upstreams {
		hosts, err := utils.GetHostsForUpstream(us)
		if err != nil {
			return nil, errors.Wrapf(err, "getting host for upstream")
		}
		for _, host := range hosts {
			uniqueHostMap[host] = struct{}{}
		}

	}
	for host := range uniqueHostMap {
		uniqueHosts = append(uniqueHosts, host)
	}
	sort.Strings(uniqueHosts)
	return uniqueHosts, nil
}

func getHostForUpstream(us *gloov1.Upstream) (string, error) {
	hosts, err := utils.GetHostsForUpstream(us)
	if err != nil {
		return "", err
	}
	if len(hosts) != 2 {
		return "", errors.Errorf("%v invalid length, expected 2", hosts)
	}
	return hosts[0], nil
}

func virtualNodesFromUpstreams(meshName string, upstreams gloov1.UpstreamList) ([]appmesh.CreateVirtualNodeInput, error) {
	portsByHost := make(map[string][]uint32)
	// TODO: filter hosts by policy, i.e. only what the user wants
	for _, us := range upstreams {
		host, err := getHostForUpstream(us)
		if err != nil {
			return nil, errors.Wrapf(err, "getting host for upstream")
		}
		port, err := utils.GetPortForUpstream(us)
		if err != nil {
			return nil, errors.Wrapf(err, "getting port for upstream")
		}
		// only add unique ports
		var alreadyAdded bool
		for _, addedPort := range portsByHost[host] {
			if addedPort == port {
				alreadyAdded = true
				break
			}
		}
		if !alreadyAdded {
			portsByHost[host] = append(portsByHost[host], port)
		}
	}
	uniqueHosts, err := getUniqueHosts(upstreams)
	if err != nil {
		return nil, errors.Wrapf(err, "getting unique hosts")
	}

	var virtualNodes []appmesh.CreateVirtualNodeInput
	for host, ports := range portsByHost {
		var listeners []*appmesh.Listener
		for _, port := range ports {
			listener := &appmesh.Listener{
				PortMapping: &appmesh.PortMapping{
					// TODO: support more than just http here
					Protocol: aws.String("http"),
					Port:     aws.Int64(int64(port)),
				},
			}
			listeners = append(listeners, listener)
		}
		virtualNode := appmesh.CreateVirtualNodeInput{
			MeshName: aws.String(meshName),
			Spec: &appmesh.VirtualNodeSpec{
				Backends:  aws.StringSlice(uniqueHosts),
				Listeners: listeners,
				ServiceDiscovery: &appmesh.ServiceDiscovery{
					Dns: &appmesh.DnsServiceDiscovery{
						ServiceName: aws.String(host),
					},
				},
			},
			VirtualNodeName: aws.String(nameutils.SanitizeName(host)),
		}
		virtualNodes = append(virtualNodes, virtualNode)
	}
	sort.SliceStable(virtualNodes, func(i, j int) bool {
		return virtualNodes[i].String() < virtualNodes[j].String()
	})
	return virtualNodes, nil
}

func virtualRoutersFromUpstreams(meshName string, upstreams gloov1.UpstreamList) ([]appmesh.CreateVirtualRouterInput, error) {
	var virtualNodes []appmesh.CreateVirtualRouterInput
	uniqueHosts, err := getUniqueHosts(upstreams)
	if err != nil {
		return nil, errors.Wrapf(err, "getting unique hosts for virtual router")
	}
	for _, host := range uniqueHosts {
		virtualNode := appmesh.CreateVirtualRouterInput{
			MeshName:          aws.String(meshName),
			VirtualRouterName: aws.String(nameutils.SanitizeName(host)),
			Spec: &appmesh.VirtualRouterSpec{
				ServiceNames: aws.StringSlice([]string{host}),
			},
		}
		virtualNodes = append(virtualNodes, virtualNode)
	}
	sort.SliceStable(virtualNodes, func(i, j int) bool {
		return virtualNodes[i].String() < virtualNodes[j].String()
	})
	return virtualNodes, nil
}

func rulesByHost(upstreams gloov1.UpstreamList, routingRules v1.RoutingRuleList) (map[string]v1.RoutingRuleList, error) {
	rules := make(map[string]v1.RoutingRuleList)
	for _, us := range upstreams {
		host, err := getHostForUpstream(us)
		if err != nil {
			return nil, errors.Wrapf(err, "getting host for upstream")
		}
		rules[host] = v1.RoutingRuleList{}
	}
	// add copy each rule for the hosts it applies to
	for _, rule := range routingRules {
		for _, usRef := range rule.Destinations {
			us, err := upstreams.Find(usRef.Strings())
			if err != nil {
				return nil, errors.Wrapf(err, "cannot find destination for routing rule")
			}
			host, err := getHostForUpstream(us)
			if err != nil {
				return nil, errors.Wrapf(err, "getting host for upstream")
			}
			rules[host] = append(rules[host], rule)
		}
	}
	return rules, nil
}

func routesForUpstreams(meshName string, upstreams gloov1.UpstreamList, routingRules v1.RoutingRuleList) ([]appmesh.CreateRouteInput, error) {
	// todo: using selector, figure out which source upstreams and which destinations need
	// the route. we are only going to support traffic shifting for now
	var routes []appmesh.CreateRouteInput

	hostsAndRules, err := rulesByHost(upstreams, routingRules)
	if err != nil {
		return nil, errors.Wrapf(err, "applying routing rules by host")
	}

	for host, routingRules := range hostsAndRules {
		var targets []*appmesh.WeightedTarget
		prefix := "/"
		// currently only 1 or 0 route rules are supported
		// aws may increase supported limits in the future
		if len(routingRules) > 0 {
			if len(routingRules) > 1 {
				return nil, errors.Errorf("currently only one routing rule is supporter per host")
			}
			rule := routingRules[0]

			if len(rule.RequestMatchers) > 0 {
				// TODO: when appmesh supports multiple matchers, we should too
				// for now, just pick the first one
				match := rule.RequestMatchers[0]
				// TODO: when appmesh supports more types of path matching, we should too
				if prefixSpecifier, ok := match.PathSpecifier.(*gloov1.Matcher_Prefix); ok {
					prefix = prefixSpecifier.Prefix
				}
			}
			// NOTE: only 3 destinations are allowed at time of release
			if rule.TrafficShifting != nil && len(rule.TrafficShifting.Destinations) > 0 {
				for _, dest := range rule.TrafficShifting.Destinations {
					us, err := upstreams.Find(dest.Upstream.Strings())
					if err != nil {
						return nil, errors.Wrapf(err, "cannot find destination for routing rule")
					}
					destinationHost, err := getHostForUpstream(us)
					if err != nil {
						return nil, errors.Wrapf(err, "getting host for destination upstream")
					}
					targets = append(targets, &appmesh.WeightedTarget{
						VirtualNode: aws.String(nameutils.SanitizeName(destinationHost)),
						Weight:      aws.Int64(int64(dest.Weight)),
					})
				}
			}
		}

		if len(targets) == 0 {
			// this is a default route, just point it at destination host
			targets = append(targets, &appmesh.WeightedTarget{
				VirtualNode: aws.String(nameutils.SanitizeName(host)),
				Weight:      aws.Int64(int64(100)),
			})
		}
		route := appmesh.CreateRouteInput{
			MeshName: aws.String(meshName),
			// TODO: names will have to be made more unique when route limits increase
			RouteName:         aws.String(nameutils.SanitizeName(host)),
			VirtualRouterName: aws.String(nameutils.SanitizeName(host)),
			Spec: &appmesh.RouteSpec{
				HttpRoute: &appmesh.HttpRoute{
					Match: &appmesh.HttpRouteMatch{
						Prefix: aws.String(prefix),
					},
					Action: &appmesh.HttpRouteAction{
						WeightedTargets: targets,
					},
				},
			},
		}
		routes = append(routes, route)
	}

	sort.SliceStable(routes, func(i, j int) bool {
		return routes[i].String() < routes[j].String()
	})
	return routes, nil
}

func purgeMesh(client AppMeshClient, mesh *appmesh.MeshRef) error {
	// delete all vrouters & routes for mesh
	virtualRouters, err := client.ListVirtualRouters(&appmesh.ListVirtualRoutersInput{
		MeshName: mesh.MeshName,
	})
	if err != nil {
		return errors.Wrapf(err, "listing existing meshes")
	}
	for _, vRouter := range virtualRouters.VirtualRouters {
		routesForRouter, err := client.ListRoutes(&appmesh.ListRoutesInput{
			MeshName:          vRouter.MeshName,
			VirtualRouterName: vRouter.VirtualRouterName,
		})
		if err != nil {
			return errors.Wrapf(err, "listing routes for %v", *vRouter.VirtualRouterName)
		}
		for _, route := range routesForRouter.Routes {
			if _, err := client.DeleteRoute(&appmesh.DeleteRouteInput{
				MeshName:          vRouter.MeshName,
				VirtualRouterName: vRouter.VirtualRouterName,
				RouteName:         route.RouteName,
			}); err != nil {
				return errors.Wrapf(err, "deleting route %v for %v", *route.RouteName, *vRouter.VirtualRouterName)
			}
		}
		if _, err := client.DeleteVirtualRouter(&appmesh.DeleteVirtualRouterInput{
			MeshName:          vRouter.MeshName,
			VirtualRouterName: vRouter.VirtualRouterName,
		}); err != nil {
			return errors.Wrapf(err, "deleting stale resource %v", vRouter)
		}
	}
	// delete virtual nodes for mesh
	existingVirtualNodes, err := client.ListVirtualNodes(&appmesh.ListVirtualNodesInput{
		MeshName: mesh.MeshName,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to list existing virtual nodes for mesh %v", *mesh.MeshName)
	}
	for _, original := range existingVirtualNodes.VirtualNodes {
		if _, err := client.DeleteVirtualNode(&appmesh.DeleteVirtualNodeInput{
			MeshName:        original.MeshName,
			VirtualNodeName: original.VirtualNodeName,
		}); err != nil {
			return errors.Wrapf(err, "deleting mesh dependency vnode %v", original)
		}
	}
	// delete mesh itself
	if _, err := client.DeleteMesh(&appmesh.DeleteMeshInput{
		MeshName: mesh.MeshName,
	}); err != nil {
		return errors.Wrapf(err, "deleting mesh %v", *mesh.MeshName)
	}
	return nil
}

/*
reconciliation
*/
func reconcileMesh(client AppMeshClient, mesh appmesh.CreateMeshInput) error {
	_, err := client.DescribeMesh(&appmesh.DescribeMeshInput{
		MeshName: mesh.MeshName,
	})
	if err == nil {
		return nil
	}
	aerr, ok := err.(awserr.Error)
	if !ok || aerr.Code() != appmesh.ErrCodeNotFoundException {
		return errors.Wrapf(err, "describing mesh %v", mesh.MeshName)
	}
	// we got a not found error, time to clean up existing meshes before installing our new one
	// TODO: when service limits go up on AWS, maybe increase the number of allowed meshes

	existingMeshes, err := client.ListMeshes(&appmesh.ListMeshesInput{})
	if err != nil {
		return errors.Wrapf(err, "listing existing meshes")
	}
	if len(existingMeshes.Meshes) > 0 {
		// clean up all previous mesh
		for _, existingMesh := range existingMeshes.Meshes {
			if err := purgeMesh(client, existingMesh); err != nil {
				return errors.Wrapf(err, "purging stale mesh for reconcile")
			}
		}
	}
	if _, err := client.CreateMesh(&mesh); err != nil {
		return errors.Wrapf(err, "creating new mesh")
	}
	return nil
}

func reconcileVirtualNodes(client AppMeshClient, mesh appmesh.CreateMeshInput, vNodes []appmesh.CreateVirtualNodeInput) error {
	existingVirtualNodes, err := client.ListVirtualNodes(&appmesh.ListVirtualNodesInput{
		MeshName: mesh.MeshName,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to list existing virtual nodes for mesh %v", *mesh.MeshName)
	}
	// delete unused
	for _, original := range existingVirtualNodes.VirtualNodes {
		cleanup := true
		for _, desired := range vNodes {
			if *desired.VirtualNodeName == aws.StringValue(original.VirtualNodeName) {
				cleanup = false
				break
			}
		}
		if cleanup {
			if _, err := client.DeleteVirtualNode(&appmesh.DeleteVirtualNodeInput{
				MeshName:        original.MeshName,
				VirtualNodeName: original.VirtualNodeName,
			}); err != nil {
				return errors.Wrapf(err, "deleting stale resource %v", original)
			}
		}
	}
	for _, desiredVNode := range vNodes {
		if err := reconcileVirtualNode(client, desiredVNode, existingVirtualNodes.VirtualNodes); err != nil {
			return errors.Wrapf(err, "reconciling virtual node %v", *desiredVNode.VirtualNodeName)
		}
	}
	return nil
}

func reconcileVirtualNode(client AppMeshClient, desiredVNode appmesh.CreateVirtualNodeInput, existingVirtualNodes []*appmesh.VirtualNodeRef) error {
	for _, node := range existingVirtualNodes {
		if aws.StringValue(node.VirtualNodeName) != *desiredVNode.VirtualNodeName {
			continue
		}
		// update
		originalVNode, err := client.DescribeVirtualNode(&appmesh.DescribeVirtualNodeInput{
			MeshName:        desiredVNode.MeshName,
			VirtualNodeName: desiredVNode.VirtualNodeName,
		})
		if err != nil {
			return errors.Wrapf(err, "retrieving original node for update")
		}
		// TODO: find a better way of comparing AWS structs
		if originalVNode.VirtualNode.Spec.String() == desiredVNode.Spec.String() {
			// spec already matches, nothing to do
			return nil
		}
		if _, err := client.UpdateVirtualNode(&appmesh.UpdateVirtualNodeInput{
			MeshName:        desiredVNode.MeshName,
			VirtualNodeName: desiredVNode.VirtualNodeName,
			Spec:            desiredVNode.Spec,
		}); err != nil {
			return errors.Wrapf(err, "updating virtual node")
		}
		return nil
	}
	// create new
	if _, err := client.CreateVirtualNode(&desiredVNode); err != nil {
		return errors.Wrapf(err, "creating virtual node")
	}

	return nil
}

func reconcileVirtualRouters(client AppMeshClient, mesh appmesh.CreateMeshInput, vRouters []appmesh.CreateVirtualRouterInput) error {
	existingVirtualRouters, err := client.ListVirtualRouters(&appmesh.ListVirtualRoutersInput{
		MeshName: mesh.MeshName,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to list existing virtual routers for mesh %v", *mesh.MeshName)
	}
	// delete unused
	for _, original := range existingVirtualRouters.VirtualRouters {
		cleanup := true
		for _, desired := range vRouters {
			if *desired.VirtualRouterName == aws.StringValue(original.VirtualRouterName) {
				cleanup = false
				break
			}
		}
		if cleanup {
			routesForRouter, err := client.ListRoutes(&appmesh.ListRoutesInput{
				MeshName:          original.MeshName,
				VirtualRouterName: original.VirtualRouterName,
			})
			if err != nil {
				return errors.Wrapf(err, "listing routes for %v", *original.VirtualRouterName)
			}
			for _, route := range routesForRouter.Routes {
				if _, err := client.DeleteRoute(&appmesh.DeleteRouteInput{
					MeshName:          original.MeshName,
					VirtualRouterName: original.VirtualRouterName,
					RouteName:         route.RouteName,
				}); err != nil {
					return errors.Wrapf(err, "deleting route %v for %v", *route.RouteName, *original.VirtualRouterName)
				}
			}
			if _, err := client.DeleteVirtualRouter(&appmesh.DeleteVirtualRouterInput{
				MeshName:          original.MeshName,
				VirtualRouterName: original.VirtualRouterName,
			}); err != nil {
				return errors.Wrapf(err, "deleting stale resource %v", original)
			}
		}
	}
	for _, desiredVRouter := range vRouters {
		if err := reconcileVirtualRouter(client, desiredVRouter, existingVirtualRouters.VirtualRouters); err != nil {
			return errors.Wrapf(err, "reconciling virtual router %v", *desiredVRouter.VirtualRouterName)
		}
	}
	return nil
}

func reconcileVirtualRouter(client AppMeshClient, desiredVRouter appmesh.CreateVirtualRouterInput, existingVirtualRouters []*appmesh.VirtualRouterRef) error {
	for _, router := range existingVirtualRouters {
		if aws.StringValue(router.VirtualRouterName) == *desiredVRouter.VirtualRouterName {
			// update
			originalVRouter, err := client.DescribeVirtualRouter(&appmesh.DescribeVirtualRouterInput{
				MeshName:          desiredVRouter.MeshName,
				VirtualRouterName: desiredVRouter.VirtualRouterName,
			})
			if err != nil {
				return errors.Wrapf(err, "retrieving original router for update")
			}
			// TODO: find a better way of comparing AWS structs
			if originalVRouter.VirtualRouter.Spec.String() == desiredVRouter.Spec.String() {
				// spec already matches, nothing to do
				return nil
			}
			if _, err := client.UpdateVirtualRouter(&appmesh.UpdateVirtualRouterInput{
				MeshName:          desiredVRouter.MeshName,
				VirtualRouterName: desiredVRouter.VirtualRouterName,
				Spec:              desiredVRouter.Spec,
			}); err != nil {
				return errors.Wrapf(err, "updating virtual router")
			}

			return nil
		}
	}
	if _, err := client.CreateVirtualRouter(&desiredVRouter); err != nil {
		return errors.Wrapf(err, "creating virtual router")
	}

	return nil
}

func reconcileRoutes(client AppMeshClient, mesh appmesh.CreateMeshInput, vRouterName string, routes []appmesh.CreateRouteInput) error {
	existingRoutes, err := client.ListRoutes(&appmesh.ListRoutesInput{
		MeshName:          mesh.MeshName,
		VirtualRouterName: aws.String(vRouterName),
	})
	if err != nil {
		return errors.Wrapf(err, "failed to list existing routes for mesh %v", *mesh.MeshName)
	}
	// delete unused
	for _, original := range existingRoutes.Routes {
		if aws.StringValue(original.VirtualRouterName) != vRouterName {
			continue
		}
		cleanup := true
		for _, desired := range routes {
			if *desired.RouteName == aws.StringValue(original.RouteName) {
				cleanup = false
				break
			}
		}
		if cleanup {
			if _, err := client.DeleteRoute(&appmesh.DeleteRouteInput{
				MeshName:          original.MeshName,
				VirtualRouterName: original.VirtualRouterName,
				RouteName:         original.RouteName,
			}); err != nil {
				return errors.Wrapf(err, "deleting stale resource %v", original)
			}
		}
	}
	for _, desiredRoute := range routes {
		if *desiredRoute.VirtualRouterName != vRouterName {
			continue
		}
		if err := reconcileRoute(client, desiredRoute, existingRoutes.Routes); err != nil {
			return errors.Wrapf(err, "reconciling route %v", *desiredRoute.RouteName)
		}
	}
	return nil
}

func reconcileRoute(client AppMeshClient, desiredRoute appmesh.CreateRouteInput, existingRoutes []*appmesh.RouteRef) error {
	for _, router := range existingRoutes {
		if aws.StringValue(router.RouteName) == *desiredRoute.RouteName {
			// update
			originalRoute, err := client.DescribeRoute(&appmesh.DescribeRouteInput{
				MeshName:          desiredRoute.MeshName,
				VirtualRouterName: desiredRoute.VirtualRouterName,
				RouteName:         desiredRoute.RouteName,
			})
			if err != nil {
				return errors.Wrapf(err, "retrieving original route for update")
			}
			// TODO: find a better way of comparing AWS structs
			if originalRoute.Route.Spec.String() == desiredRoute.Spec.String() {
				// spec already matches, nothing to do
				return nil
			}
			if _, err := client.UpdateRoute(&appmesh.UpdateRouteInput{
				MeshName:  desiredRoute.MeshName,
				RouteName: desiredRoute.RouteName,
				Spec:      desiredRoute.Spec,
			}); err != nil {
				return errors.Wrapf(err, "updating route")
			}
			return nil
		}
	}
	if _, err := client.CreateRoute(&desiredRoute); err != nil {
		return errors.Wrapf(err, "creating route")
	}

	return nil
}
