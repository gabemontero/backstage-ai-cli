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
$ bkstg-ai new-model <kserve|kubeflow> <owner> <lifecycle> <args...>

# Access the Backstage Catalog for Entities related to AI Models
$ bkstg-ai get [location|components|resources|apis] [args...]

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
		Use:     "new-model",
		Long:    "new-model accesses one of the supported backends and builds Backstage Catalog Entity YAML with available Model metadata",
		Aliases: []string{"create", "c", "nm", "new-models"},
		Example: `
# Access a supported backend for AI Model metadata and generate Backstage Catalog Entity YAML for that metadata
`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	newModel.AddCommand(kserve.NewCmd(cfg))
	newModel.AddCommand(kubeflowmodelregistry.NewCmd(cfg))

	queryModel := &cobra.Command{
		Use:     "get",
		Long:    "get accesses the Backstage Catalog for Entities related to AI Models",
		Aliases: []string{"g"},
		Example: `
# Access the Backstage Catalog for Entities related to AI Models
$ bkstg-ai get <locations|components|resources|apis|entities> [args...]
`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	deleteModel := &cobra.Command{
		Use:     "delete-model",
		Long:    "delete-model removes the Backstage Catalog for Entities corresponding to the provided location ID",
		Aliases: []string{"delete", "dm", "del", "d", "delete-models"},
		Example: `
# Remove from the Backstage Catalog the Location entity for the provided Location ID, using the dynamically generated 
# hash ID from when the location was imported.  There is not support in Backstage currently for specifying
# the URL used to import the model as a query parameter.
$ bkstg-ai delete-model <location id>

# Set the URL for the Backstage instance, the authentication token, and Skip-TLS settings 
$ bkstg-ai delete-model <location id> --backstage-url=https://my-rhdh.com --backstage-token=my-token --backstage-skip-tls=true
`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				klog.Error("ERROR: delete-model requires a location ID")
			}
			processOutput(backstage.SetupBackstageRESTClient(cfg).DeleteLocation(args[0]))
		},
	}
	importModel := &cobra.Command{
		Use:     "import-model",
		Long:    "import-model updates the Backstage Catalog with Entities contained in the provided location URL",
		Aliases: []string{"post", "im", "p", "i", "import-models"},
		Example: `
# Import from an accessible URL Backstage Catalog entities
$ bkstg-ai import-model <url>

# Set the additional URL for the Backstage instance, the authentication token, and Skip-TLS settings 
$ bkstg-ai import-model <url> --backstage-url=https://my-rhdh.com --backstage-token=my-token --backstage-skip-tls=true
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
		Use:     "entities",
		Long:    "entities retrieves the AI related Backstage Catalog Entities",
		Aliases: []string{"e", "entity"},
		Example: `
# Access the Backstage Catalog for all entities, regardless if AI related
$ bkstg-ai get entities

# Set the URL for the Backstage, the authentication token, and Skip-TLS settings
$ bkstg-ai get entities --backstage-url=https://my-rhdh.com --backstage-token=my-token --backstage-skip-tls=true
`,
		RunE: func(cmd *cobra.Command, args []string) error {

			str, err := backstage.SetupBackstageRESTClient(cfg).ListEntities()
			processOutput(str, err)
			return err

		},
	})

	queryModel.AddCommand(&cobra.Command{
		Use:     "locations",
		Long:    "locations retrieves the AI related Backstage Catalog Locations",
		Aliases: []string{"l", "location"},
		Example: `
# Access the Backstage Catalog for locations, regardless if AI related
$ bkstg-ai get locations [args...]

# Access the Backstage Catatlog for a specific location using the dynamically generated 
# hash ID from when the location was imported.  There is not support in Backstage currently for specifying
# the URL used to import the model as a query parameter.
$ bkstg-ai get locations my-big-long-id-for-location

# Set the URL for the Backstage, the authentication token, and Skip-TLS settings
$ bkstg-ai get locations --backstage-url=https://my-rhdh.com --backstage-token=my-token --backstage-skip-tls=true
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			str, err := backstage.SetupBackstageRESTClient(cfg).GetLocation(args...)
			processOutput(str, err)
			return err
		},
	})

	queryModel.AddCommand(&cobra.Command{
		Use:     "components",
		Long:    "components retrieves the AI related Backstage Catalog Components",
		Aliases: []string{"c", "component"},
		Example: `
# Retrieve the Backstage Catalog for resources related to AI Models, where being AI related is determined by the 
# 'type' being set to 'model-server'
$ bkstg-ai get components [args...]

# Set the URL for the Backstage, the authentication token, and Skip-TLS settings
$ bkstg-ai get components --backstage-url=https://my-rhdh.com --backstage-token=my-token --backstage-skip-tls=true

# Retrieve a specific set of AI related Components by namespace:name
$ bkstg-ai get components default:my-component default:your-component

# Retrieve a set of AI Components where the provided list of tags match (order of tags disregarded)
$ bkstg-ai get components genai vllm --use-params-as-tags=true

# Retrieve a set of Components which have any of the provided list of tags
$ bkstg-ai get components gen-ai --use-params-as-tags=true --use-any-subset=true
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			str, err := backstage.SetupBackstageRESTClient(cfg).GetComponent(args...)
			processOutput(str, err)
			return err
		},
	})

	queryModel.AddCommand(&cobra.Command{
		Use:     "resources",
		Long:    "resources retrieves the AI related Backstage Catalog Resources",
		Aliases: []string{"r", "resource"},
		Example: `
# Retrieve the Backstage Catalog for resources related to AI Models, where being AI related is determined by the 
# 'type' being set to 'ai-model'
$ bkstg-ai get resources [args...]

# Set the URL for the Backstage, the authentication token, and Skip-TLS settings
$ bkstg-ai get resources --backstage-url=https://my-rhdh.com --backstage-token=my-token --backstage-skip-tls=true

# Retrieve a specific set of AI related Resources by namespace:name
$ bkstg-ai get resources default:my-component default:your-component

# Retrieve a set of AI Resources where the provided list of tags match (order of tags disregarded)
$ bkstg-ai get resources genai vllm --use-params-as-tags=true

# Retrieve a set of AI Resources which have any of the provided list of tags
$ bkstg-ai get resources gen-ai --use-params-as-tags=true --use-any-subset=true
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			str, err := backstage.SetupBackstageRESTClient(cfg).GetResource(args...)
			processOutput(str, err)
			return err
		},
	})

	queryModel.AddCommand(&cobra.Command{
		Use:     "apis",
		Long:    "apis retrieves the AI related Backstage Catalog APIS",
		Aliases: []string{"a", "api"},
		Example: `
# Retrieve the Backstage Catalog for APIs related to AI Models, where being AI related is determined by the 
# 'type' being set to 'model-service-api'
$ bkstg-ai get apis [args...]

# Set the URL for the Backstage, the authentication token, and Skip-TLS settings
$ bkstg-ai get locations --backstage-url=https://my-rhdh.com --backstage-token=my-token --backstage-skip-tls=true

# Retrieve a specific set of AI related APIs by namespace:name
$ bkstg-ai get apis default:my-component default:your-component

# Retrieve a set of AI APIs where the provided list of tags match (order of tags disregarded)
$ bkstg-ai get apis genai vllm --use-params-as-tags=true

# Retrieve a set of AI APIs which have any of the provided list of tags
$ bkstg-ai get apis gen-ai --use-params-as-tags=true --use-any-subset=true
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			str, err := backstage.SetupBackstageRESTClient(cfg).GetAPI(args...)
			processOutput(str, err)
			return err
		},
	})

	queryModel.PersistentFlags().BoolVar(&(cfg.ParamsAsTags), "use-params-as-tags", cfg.ParamsAsTags,
		"Use any additional parameters as tag identifiers")
	queryModel.PersistentFlags().BoolVar(&(cfg.AnySubsetWorks), "use-any-subset", cfg.AnySubsetWorks,
		"Use any additional parameters as tag identifiers")

	return bkstgAI
}

func processOutput(str string, err error) {
	klog.Infoln(str)
	klog.Flush()
	if err != nil {
		klog.Errorf("%s", err.Error())
		klog.Flush()
	}
}
