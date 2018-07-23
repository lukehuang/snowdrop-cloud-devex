package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/openshift/api/apps/v1"
	"github.com/snowdrop/k8s-supervisor/pkg/buildpack"
	"github.com/snowdrop/k8s-supervisor/pkg/buildpack/types"
)

var initCmd = &cobra.Command{
	Use:     "init [flags]",
	Short:   "Create a development's pod for the component",
	Long:    `Create a development's pod for the component.`,
	Example: ` sb init -n bootapp`,
	Args:    cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {

		log.Info("Init command called")
		setup := Setup()

		// Create ImageStreams
		log.Info("Create ImageStreams for Supervisord and Java S2I Image of SpringBoot")
		images := []types.Image{
			*buildpack.CreateTypeImage(true, "dev-s2i", "latest", "quay.io/snowdrop/spring-boot-s2i", false),
			*buildpack.CreateTypeImage(true, "copy-supervisord", "latest", "quay.io/snowdrop/supervisord", true),
		}
		buildpack.CreateImageStreamTemplate(setup.RestConfig, setup.Application, images)

		// Create PVC
		log.Info("Create PVC to store m2 repo")
		buildpack.CreatePVC(setup.Clientset, setup.Application, "1Gi")

		var dc *v1.DeploymentConfig
		log.Info("Create or retrieve DeploymentConfig using Supervisord and Java S2I Image of SpringBoot")
		dc = buildpack.CreateOrRetrieveDeploymentConfig(setup.RestConfig, setup.Application)

		log.Info("Create Service using Template")
		buildpack.CreateServiceTemplate(setup.Clientset, dc, setup.Application)

		log.Info("Create Route using Template")
		buildpack.CreateRouteTemplate(setup.RestConfig, setup.Application)
	},
}

func init() {

	// Add a defined annotation in order to appear in the help menu
	initCmd.Annotations = map[string]string{"command": "init"}

	rootCmd.AddCommand(initCmd)
}
