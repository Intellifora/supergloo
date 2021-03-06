// Code generated by solo-kit. DO NOT EDIT.

// +build solokit

package v1alpha1

import (
	"context"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/go-utils/kubeutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	kuberc "github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/utils/log"
	"github.com/solo-io/solo-kit/test/helpers"
	"github.com/solo-io/solo-kit/test/setup"
	"k8s.io/client-go/rest"

	// Needed to run tests in GKE
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	// From https://github.com/kubernetes/client-go/blob/53c7adfd0294caa142d961e1f780f74081d5b15f/examples/out-of-cluster-client-configuration/main.go#L31
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var _ = Describe("V1Alpha1Emitter", func() {
	if os.Getenv("RUN_KUBE_TESTS") != "1" {
		log.Printf("This test creates kubernetes resources and is disabled by default. To enable, set RUN_KUBE_TESTS=1 in your env.")
		return
	}
	var (
		namespace1               string
		namespace2               string
		name1, name2             = "angela" + helpers.RandString(3), "bob" + helpers.RandString(3)
		cfg                      *rest.Config
		emitter                  RbacEmitter
		serviceRoleClient        ServiceRoleClient
		serviceRoleBindingClient ServiceRoleBindingClient
		rbacConfigClient         RbacConfigClient
	)

	BeforeEach(func() {
		namespace1 = helpers.RandString(8)
		namespace2 = helpers.RandString(8)
		var err error
		cfg, err = kubeutils.GetConfig("", "")
		Expect(err).NotTo(HaveOccurred())
		err = setup.SetupKubeForTest(namespace1)
		Expect(err).NotTo(HaveOccurred())
		err = setup.SetupKubeForTest(namespace2)
		Expect(err).NotTo(HaveOccurred())
		// ServiceRole Constructor
		serviceRoleClientFactory := &factory.KubeResourceClientFactory{
			Crd:         ServiceRoleCrd,
			Cfg:         cfg,
			SharedCache: kuberc.NewKubeCache(),
		}
		serviceRoleClient, err = NewServiceRoleClient(serviceRoleClientFactory)
		Expect(err).NotTo(HaveOccurred())
		// ServiceRoleBinding Constructor
		serviceRoleBindingClientFactory := &factory.KubeResourceClientFactory{
			Crd:         ServiceRoleBindingCrd,
			Cfg:         cfg,
			SharedCache: kuberc.NewKubeCache(),
		}
		serviceRoleBindingClient, err = NewServiceRoleBindingClient(serviceRoleBindingClientFactory)
		Expect(err).NotTo(HaveOccurred())
		// RbacConfig Constructor
		rbacConfigClientFactory := &factory.KubeResourceClientFactory{
			Crd:         RbacConfigCrd,
			Cfg:         cfg,
			SharedCache: kuberc.NewKubeCache(),
		}
		rbacConfigClient, err = NewRbacConfigClient(rbacConfigClientFactory)
		Expect(err).NotTo(HaveOccurred())
		emitter = NewRbacEmitter(serviceRoleClient, serviceRoleBindingClient, rbacConfigClient)
	})
	AfterEach(func() {
		setup.TeardownKube(namespace1)
		setup.TeardownKube(namespace2)
	})
	It("tracks snapshots on changes to any resource", func() {
		ctx := context.Background()
		err := emitter.Register()
		Expect(err).NotTo(HaveOccurred())

		snapshots, errs, err := emitter.Snapshots([]string{namespace1, namespace2}, clients.WatchOpts{
			Ctx:         ctx,
			RefreshRate: time.Second,
		})
		Expect(err).NotTo(HaveOccurred())

		var snap *RbacSnapshot

		/*
			ServiceRole
		*/

		assertSnapshotServiceRoles := func(expectServiceRoles ServiceRoleList, unexpectServiceRoles ServiceRoleList) {
		drain:
			for {
				select {
				case snap = <-snapshots:
					for _, expected := range expectServiceRoles {
						if _, err := snap.ServiceRoles.List().Find(expected.Metadata.Ref().Strings()); err != nil {
							continue drain
						}
					}
					for _, unexpected := range unexpectServiceRoles {
						if _, err := snap.ServiceRoles.List().Find(unexpected.Metadata.Ref().Strings()); err == nil {
							continue drain
						}
					}
					break drain
				case err := <-errs:
					Expect(err).NotTo(HaveOccurred())
				case <-time.After(time.Second * 10):
					nsList1, _ := serviceRoleClient.List(namespace1, clients.ListOpts{})
					nsList2, _ := serviceRoleClient.List(namespace2, clients.ListOpts{})
					combined := nsList1.ByNamespace()
					combined.Add(nsList2...)
					Fail("expected final snapshot before 10 seconds. expected " + log.Sprintf("%v", combined))
				}
			}
		}
		serviceRole1a, err := serviceRoleClient.Write(NewServiceRole(namespace1, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		serviceRole1b, err := serviceRoleClient.Write(NewServiceRole(namespace2, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotServiceRoles(ServiceRoleList{serviceRole1a, serviceRole1b}, nil)
		serviceRole2a, err := serviceRoleClient.Write(NewServiceRole(namespace1, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		serviceRole2b, err := serviceRoleClient.Write(NewServiceRole(namespace2, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotServiceRoles(ServiceRoleList{serviceRole1a, serviceRole1b, serviceRole2a, serviceRole2b}, nil)

		err = serviceRoleClient.Delete(serviceRole2a.Metadata.Namespace, serviceRole2a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = serviceRoleClient.Delete(serviceRole2b.Metadata.Namespace, serviceRole2b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotServiceRoles(ServiceRoleList{serviceRole1a, serviceRole1b}, ServiceRoleList{serviceRole2a, serviceRole2b})

		err = serviceRoleClient.Delete(serviceRole1a.Metadata.Namespace, serviceRole1a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = serviceRoleClient.Delete(serviceRole1b.Metadata.Namespace, serviceRole1b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotServiceRoles(nil, ServiceRoleList{serviceRole1a, serviceRole1b, serviceRole2a, serviceRole2b})

		/*
			ServiceRoleBinding
		*/

		assertSnapshotServiceRoleBindings := func(expectServiceRoleBindings ServiceRoleBindingList, unexpectServiceRoleBindings ServiceRoleBindingList) {
		drain:
			for {
				select {
				case snap = <-snapshots:
					for _, expected := range expectServiceRoleBindings {
						if _, err := snap.ServiceRoleBindings.List().Find(expected.Metadata.Ref().Strings()); err != nil {
							continue drain
						}
					}
					for _, unexpected := range unexpectServiceRoleBindings {
						if _, err := snap.ServiceRoleBindings.List().Find(unexpected.Metadata.Ref().Strings()); err == nil {
							continue drain
						}
					}
					break drain
				case err := <-errs:
					Expect(err).NotTo(HaveOccurred())
				case <-time.After(time.Second * 10):
					nsList1, _ := serviceRoleBindingClient.List(namespace1, clients.ListOpts{})
					nsList2, _ := serviceRoleBindingClient.List(namespace2, clients.ListOpts{})
					combined := nsList1.ByNamespace()
					combined.Add(nsList2...)
					Fail("expected final snapshot before 10 seconds. expected " + log.Sprintf("%v", combined))
				}
			}
		}
		serviceRoleBinding1a, err := serviceRoleBindingClient.Write(NewServiceRoleBinding(namespace1, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		serviceRoleBinding1b, err := serviceRoleBindingClient.Write(NewServiceRoleBinding(namespace2, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotServiceRoleBindings(ServiceRoleBindingList{serviceRoleBinding1a, serviceRoleBinding1b}, nil)
		serviceRoleBinding2a, err := serviceRoleBindingClient.Write(NewServiceRoleBinding(namespace1, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		serviceRoleBinding2b, err := serviceRoleBindingClient.Write(NewServiceRoleBinding(namespace2, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotServiceRoleBindings(ServiceRoleBindingList{serviceRoleBinding1a, serviceRoleBinding1b, serviceRoleBinding2a, serviceRoleBinding2b}, nil)

		err = serviceRoleBindingClient.Delete(serviceRoleBinding2a.Metadata.Namespace, serviceRoleBinding2a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = serviceRoleBindingClient.Delete(serviceRoleBinding2b.Metadata.Namespace, serviceRoleBinding2b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotServiceRoleBindings(ServiceRoleBindingList{serviceRoleBinding1a, serviceRoleBinding1b}, ServiceRoleBindingList{serviceRoleBinding2a, serviceRoleBinding2b})

		err = serviceRoleBindingClient.Delete(serviceRoleBinding1a.Metadata.Namespace, serviceRoleBinding1a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = serviceRoleBindingClient.Delete(serviceRoleBinding1b.Metadata.Namespace, serviceRoleBinding1b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotServiceRoleBindings(nil, ServiceRoleBindingList{serviceRoleBinding1a, serviceRoleBinding1b, serviceRoleBinding2a, serviceRoleBinding2b})

		/*
			RbacConfig
		*/

		assertSnapshotRbacConfigs := func(expectRbacConfigs RbacConfigList, unexpectRbacConfigs RbacConfigList) {
		drain:
			for {
				select {
				case snap = <-snapshots:
					for _, expected := range expectRbacConfigs {
						if _, err := snap.RbacConfigs.List().Find(expected.Metadata.Ref().Strings()); err != nil {
							continue drain
						}
					}
					for _, unexpected := range unexpectRbacConfigs {
						if _, err := snap.RbacConfigs.List().Find(unexpected.Metadata.Ref().Strings()); err == nil {
							continue drain
						}
					}
					break drain
				case err := <-errs:
					Expect(err).NotTo(HaveOccurred())
				case <-time.After(time.Second * 10):
					nsList1, _ := rbacConfigClient.List(namespace1, clients.ListOpts{})
					nsList2, _ := rbacConfigClient.List(namespace2, clients.ListOpts{})
					combined := nsList1.ByNamespace()
					combined.Add(nsList2...)
					Fail("expected final snapshot before 10 seconds. expected " + log.Sprintf("%v", combined))
				}
			}
		}
		rbacConfig1a, err := rbacConfigClient.Write(NewRbacConfig(namespace1, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		rbacConfig1b, err := rbacConfigClient.Write(NewRbacConfig(namespace2, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotRbacConfigs(RbacConfigList{rbacConfig1a, rbacConfig1b}, nil)
		rbacConfig2a, err := rbacConfigClient.Write(NewRbacConfig(namespace1, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		rbacConfig2b, err := rbacConfigClient.Write(NewRbacConfig(namespace2, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotRbacConfigs(RbacConfigList{rbacConfig1a, rbacConfig1b, rbacConfig2a, rbacConfig2b}, nil)

		err = rbacConfigClient.Delete(rbacConfig2a.Metadata.Namespace, rbacConfig2a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = rbacConfigClient.Delete(rbacConfig2b.Metadata.Namespace, rbacConfig2b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotRbacConfigs(RbacConfigList{rbacConfig1a, rbacConfig1b}, RbacConfigList{rbacConfig2a, rbacConfig2b})

		err = rbacConfigClient.Delete(rbacConfig1a.Metadata.Namespace, rbacConfig1a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = rbacConfigClient.Delete(rbacConfig1b.Metadata.Namespace, rbacConfig1b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotRbacConfigs(nil, RbacConfigList{rbacConfig1a, rbacConfig1b, rbacConfig2a, rbacConfig2b})
	})
})
