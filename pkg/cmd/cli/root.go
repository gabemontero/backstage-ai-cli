package cli

import (
	"github.com/gabemontero/backstage-ai-cli/pkg/cmd/cli/backstage"
	"github.com/gabemontero/backstage-ai-cli/pkg/cmd/cli/kserve"
	"github.com/gabemontero/backstage-ai-cli/pkg/cmd/cli/kubeflowmodelregistry"
	"github.com/gabemontero/backstage-ai-cli/pkg/config"
	"github.com/gabemontero/backstage-ai-cli/pkg/util"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
	"os"
	"strconv"
)

// NewCmd create a new root command, linking together all sub-commands organized by groups.
func NewCmd() *cobra.Command {
	cfg := &config.Config{}
	bkstgAI := &cobra.Command{
		Use:  "bkstg-ai",
		Long: "Backstage AI is a command line tool that facilitates management of AI related Entities in the Backstage Catalog.",
		Example: `
# Access a supported backend for AI Model metadata and generate Backstage Catalog Entity YAML for that metadata
$ bkstg-ai new-model <kserve|kubeflow|huggingface|oci|3scale> <owner> <lifecycle> <args...>

# Access the Backstage Catalog for Entities related to AI Models
$ bkstg-ai fetch-model [with-any-tags|with-all-tags] [location|components|resources|apis] [args...]

# Import from an accessible URL Backstage Catalog entities
$ bkstg-ai import-model <url>

# Remove from the Backstage Catalog the Location entity for the provided Location ID.
$ bkstg-ai delete-model <location id>
`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cfg.Kubeconfig = os.Getenv("KUBECONFIG")
	cfg.BackstageURL = os.Getenv("BACKSTAGE_URL")
	cfg.BackstageToken = os.Getenv("BACKSTAGE_TOKEN")
	cfg.BackstageSkipTLS, _ = strconv.ParseBool(os.Getenv("BACKSTAGE_SKIP_TLS"))
	cfg.StoreURL = os.Getenv("MODEL_METADATA_URL")
	cfg.StoreToken = os.Getenv("MODEL_METADATA_TOKEN")
	cfg.StoreSkipTLS, _ = strconv.ParseBool(os.Getenv("METADATA_MODEL_SKIP_TLS"))
	cfg.Namespace = util.GetCurrentProject()

	bkstgAI.PersistentFlags().StringVar(&(cfg.Kubeconfig), "kubeconfig", cfg.Kubeconfig,
		"Path to the kubeconfig file to use for CLI requests.")
	bkstgAI.PersistentFlags().StringVar(&(cfg.Namespace), "namespace", cfg.Namespace,
		"The name of the Kubernetes namespace to use for CLI requests.")
	bkstgAI.PersistentFlags().StringVar(&(cfg.BackstageURL), "backstage-url", cfg.BackstageURL,
		"The URL used for accessing the Backstage Catalog REST API.")
	bkstgAI.PersistentFlags().StringVar(&(cfg.BackstageToken), "backstage-token", cfg.BackstageToken,
		"The bearer authorization token used for accessing the Backstage Catalog REST API.")
	bkstgAI.PersistentFlags().BoolVar(&(cfg.BackstageSkipTLS), "backstage-skip-tls", cfg.StoreSkipTLS,
		"Whether to skip use of TLS when accessing the Backstage Catalog REST API.")
	bkstgAI.PersistentFlags().StringVar(&(cfg.StoreURL), "model-metadata-url", cfg.StoreURL,
		"The URL used for accessing the external source for Model Metadata.")
	bkstgAI.PersistentFlags().StringVar(&(cfg.StoreToken), "model-metadata-token", cfg.StoreToken,
		"The bearer authorization token used for accessing the external source for Model Metadata.")
	bkstgAI.PersistentFlags().BoolVar(&(cfg.StoreSkipTLS), "model-metadata-skip-tls", cfg.StoreSkipTLS,
		"Whether to skip use of TLS when accessing the external source for Model Metadata.")

	newModel := &cobra.Command{
		Use:  "new-model",
		Long: "new-model accesses one of the supported backends and builds Backstage Catalog Entity YAML with available Model metadata",
		Example: `
# Access a supported backend for AI Model metadata and generate Backstage Catalog Entity YAML for that metadata
$ bkstg-ai new-model <kserve|kubeflow|huggingface|oci|3scale> <owner> <lifecycle> <args...>
`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	newModel.AddCommand(kserve.NewCmd(cfg))
	newModel.AddCommand(kubeflowmodelregistry.NewCmd(cfg))

	queryModel := &cobra.Command{
		Use:  "fetch-model",
		Long: "fetch-model accesses the Backstage Catalog for Entities related to AI Models",
		Example: `
# Access the Backstage Catalog for Entities related to AI Models
$ bkstg-ai fetch-model [with-any-tags|with-all-tags] [location|components|resources|apis|entities] [args...]
`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	deleteModel := &cobra.Command{
		Use:  "delete-model",
		Long: "delete-model removes the Backstage Catalog for Entities corresponding to the provided location ID",
		Example: `
# Remove from the Backstage Catalog the Location entity for the provided Location ID.
$ bkstg-ai delete-model <location id>
`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				klog.Error("ERROR: delete-model requires a location ID")
			}
			processOutput(backstage.SetupBackstageRESTClient(cfg).DeleteLocation(args[0]))
		},
	}
	importModel := &cobra.Command{
		Use:  "import-model",
		Long: "import-model updates the Backstage Catalog with Entities contained in the provided location URL",
		Example: `
# Import from an accessible URL Backstage Catalog entities
$ bkstg-ai import-model <url>
`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				klog.Error("ERROR: import-model requires a location URL")
				klog.Flush()
				return
			}
			processOutput(backstage.SetupBackstageRESTClient(cfg).ImportLocation(args[0]))
		},
	}

	bkstgAI.AddCommand(newModel)
	bkstgAI.AddCommand(queryModel)
	bkstgAI.AddCommand(deleteModel)
	bkstgAI.AddCommand(importModel)

	queryModel.AddCommand(&cobra.Command{
		Use:  "entities",
		Long: "entities lists the AI related Backstage Catalog Entities",
		Example: `
# Access the Backstage Catalog for all entities, regardless if AI related
$ bkstg-ai fetch-model entities
`,
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) == 0 {
				processOutput(backstage.SetupBackstageRESTClient(cfg).ListEntities())
				return
			}
		},
	})

	queryModel.AddCommand(&cobra.Command{
		Use:  "locations",
		Long: "locations lists the AI related Backstage Catalog Locations",
		Example: `
# Access the Backstage Catalog for locations related to AI Models
$ bkstg-ai fetch-model locations [args...]
`,
		Run: func(cmd *cobra.Command, args []string) {
			processOutput(backstage.SetupBackstageRESTClient(cfg).GetLocation(args...))
		},
	})

	queryModel.AddCommand(&cobra.Command{
		Use:  "components ",
		Long: "resources retrieves the AI related Backstage Catalog Resources",
		Example: `
# Retrieve the Backstage Catalog for resources related to AI Models
$ bkstg-ai fetch-model components [args...]
`,
		Run: func(cmd *cobra.Command, args []string) {
			processOutput(backstage.SetupBackstageRESTClient(cfg).GetComponent(args...))
		},
	})

	queryModel.AddCommand(&cobra.Command{
		Use:  "resources",
		Long: "resources retrieves the AI related Backstage Catalog Resources",
		Example: `
# Retrieve the Backstage Catalog for resources related to AI Models
$ bkstg-ai fetch-model resources [args...]
`,
		Run: func(cmd *cobra.Command, args []string) {
			processOutput(backstage.SetupBackstageRESTClient(cfg).GetResource(args...))
		},
	})

	queryModel.AddCommand(&cobra.Command{
		Use:  "apis",
		Long: "apis retrieves the AI related Backstage Catalog APIS",
		Example: `
# Retrieve the Backstage Catalog for APIs related to AI Models
$ bkstg-ai fetch-model apis [args...]
`,
		Run: func(cmd *cobra.Command, args []string) {
			processOutput(backstage.SetupBackstageRESTClient(cfg).GetAPI(args...))
		},
	})

	queryModel.PersistentFlags().BoolVar(&(cfg.ParamsAsTags), "use-params-as-tags", cfg.ParamsAsTags,
		"Use any additional parameters as tag identifiers")
	queryModel.PersistentFlags().BoolVar(&(cfg.AnySubsetWorks), "use-any-subset", cfg.AnySubsetWorks,
		"Use any additional parameters as tag identifiers")

	return bkstgAI
}

func processOutput(str string, err error) {
	if err != nil {
		klog.Errorf("%s\nERROR: %s\n", str, err.Error())
		klog.Flush()
		return
	}
	klog.Infoln(str)
	klog.Flush()
}
