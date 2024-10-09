package kserve

import (
	"context"
	"github.com/gabemontero/backstage-ai-cli/pkg/config"
	"github.com/gabemontero/backstage-ai-cli/pkg/util"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

func NewCmd(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kserve",
		Short: "KServe related API",
		Long:  "Interact with KServe related instances on a K8s cluster to manage AI related catalog entities in a Backstage instance.",
		Example: `
# Import a specific inferenceservices.serving.kserve.io on a cluster as a set of Catalog, Resource, and API Entities in the Backstage catalog
$ bkstg-ai import <ID for inferenceservices.serving.kserve.io  instance>

# Import all the inferenceservices.serving.kserve.io on a cluster as a set of Catalog, Resource, and API Entities in the Backstage catalog
$ bkstg-ai import
`,
		Run: func(cmd *cobra.Command, args []string) {
			id := ""

			if len(args) > 0 {
				id = args[0]
			}

			kubeconfig, err := util.GetK8sConfig(cfg)
			if err != nil {
				klog.Errorf("ERROR: problem with kubeconfig: %v\n", err)
				return
			}
			servingClient := util.GetKServeClient(kubeconfig)
			namespace := cfg.Namespace
			if len(namespace) == 0 {
				namespace = util.GetCurrentProject()
				if len(namespace) == 0 {
					return
				}
			}
			if len(id) != 0 {
				is, err := servingClient.InferenceServices(namespace).Get(context.Background(), id, metav1.GetOptions{})
				if err != nil {
					klog.Errorf("ERROR: inference service retrieval error for %s:%s: %s\n", namespace, id, err.Error())
					return
				}
				if is.Status.URL != nil {
					//TODO start building backstage AI
				}
			}

		},
	}

	cmd.AddCommand()

	return cmd
}
