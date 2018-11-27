package consul_test

import (
	"context"

	"github.com/gogo/protobuf/types"

	"github.com/solo-io/supergloo/pkg/install"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/supergloo/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/supergloo/pkg/api/v1"
	"github.com/solo-io/supergloo/pkg/install/consul"
)

/*
Smoke test for installing and uninstalling consul
*/
var _ = Describe("Consul Installer", func() {

	installNamespace := "consul"
	superglooNamespace := "supergloo-system" // this needs to be made before running tests
	meshName := "test-consul-mesh"

	getSnapshot := func(install bool) *v1.InstallSnapshot {
		return &v1.InstallSnapshot{
			Installs: v1.InstallsByNamespace{
				superglooNamespace: v1.InstallList{
					&v1.Install{
						Metadata: core.Metadata{
							Namespace: superglooNamespace,
							Name:      meshName,
						},
						MeshType: &v1.Install_Consul{
							Consul: &v1.Consul{
								InstallationNamespace: installNamespace,
							},
						},
						ChartLocator: &v1.HelmChartLocator{
							Kind: &v1.HelmChartLocator_ChartPath{
								ChartPath: &v1.HelmChartPath{
									Path: "https://github.com/hashicorp/consul-helm/archive/v0.3.0.tar.gz",
								},
							},
						},
						Enabled: &types.BoolValue{
							Value: install,
						},
					},
				},
			},
		}
	}

	kubeCache := kube.NewKubeCache()

	var meshClient v1.MeshClient
	var syncer install.InstallSyncer

	BeforeEach(func() {
		util.TryCreateNamespace("supergloo-system")
		meshClient = util.GetMeshClient(kubeCache)
		syncer = install.InstallSyncer{
			Kube:       util.GetKubeClient(),
			MeshClient: meshClient,
		}
	})

	AfterEach(func() {
		util.TerminateNamespaceBlocking("supergloo-system")

		// just in case
		meshClient.Delete(superglooNamespace, meshName, clients.DeleteOpts{})
		util.UninstallHelmRelease(meshName)
		util.DeleteWebhookConfigIfExists(consul.WebhookCfg)
		util.DeleteCrb(consul.CrbName)
		util.TerminateNamespaceBlocking(installNamespace)
	})

	It("Can install and uninstall consul", func() {
		snap := getSnapshot(true)
		err := syncer.Sync(context.TODO(), snap)
		Expect(err).NotTo(HaveOccurred())
		Expect(util.WaitForAvailablePods(installNamespace)).To(BeEquivalentTo(2))

		snap = getSnapshot(false)
		err = syncer.Sync(context.TODO(), snap)
		Expect(err).NotTo(HaveOccurred())

		// Validate everything cleaned up
		util.WaitForTerminatedNamespace(installNamespace)
		Expect(util.HelmReleaseDoesntExist(meshName)).To(BeTrue())
		Expect(util.CrbDoesntExist(consul.CrbName)).To(BeTrue())
		Expect(util.WebhookConfigNotFound(consul.WebhookCfg)).To(BeTrue())

		mesh, err := meshClient.Read(superglooNamespace, meshName, clients.ReadOpts{})
		Expect(mesh).To(BeNil())
		Expect(err).ToNot(BeNil())
	})
})
