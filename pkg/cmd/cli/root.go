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
	"strings"
)

const (
	bkstgAIExample = `
# Access a supported backend for AI Model metadata and generate Backstage Catalog Entity YAML for that metadata
$ %s new-model <kserve|kubeflow> <owner> <lifecycle> <args...>

# Access the Backstage Catalog for Entities related to AI Models
$ %s get [location|components|resources|apis] [args...]

# Import from an accessible URL Backstage Catalog entities
$ %s import-model <url>

# Remove from the Backstage Catalog the Location entity for the provided Location ID.
$ %s delete-model <location id>
`

	newModelExample = `
# Access a supported backend for AI Model metadata and generate Backstage Catalog Entity YAML for that metadata
$ %s new-model kserve [args]
`

	getExample = `
# Access the Backstage Catalog for Entities related to AI Models
$ %s get <locations|components|resources|apis|entities> [args...]
`

	deleteModelExample = `
# Remove from the Backstage Catalog the Location entity for the provided Location ID, using the dynamically generated 
# hash ID from when the location was imported.  There is not support in Backstage currently for specifying
# the URL used to import the model as a query parameter.
$ %s delete-model <location id>

# Set the URL for the Backstage instance, the authentication token, and Skip-TLS settings 
$ %s delete-model <location id> --backstage-url=https://my-rhdh.com --backstage-token=my-token --backstage-skip-tls=true
`

	importModelExample = `
# Import from an accessible URL Backstage Catalog entities
$ %s import-model <url>

# Set the additional URL for the Backstage instance, the authentication token, and Skip-TLS settings 
$ %s import-model <url> --backstage-url=https://my-rhdh.com --backstage-token=my-token --backstage-skip-tls=true
`

	getEntitiesExample = `
# Access the Backstage Catalog for all entities, regardless if AI related
$ %s get entities

# Set the URL for the Backstage, the authentication token, and Skip-TLS settings
$ %s get entities --backstage-url=https://my-rhdh.com --backstage-token=my-token --backstage-skip-tls=true
`

	getLocationsExample = `
# Access the Backstage Catalog for locations, regardless if AI related
$ %s get locations [args...]

# Access the Backstage Catatlog for a specific location using the dynamically generated 
# hash ID from when the location was imported.  There is not support in Backstage currently for specifying
# the URL used to import the model as a query parameter.
$ %s get locations my-big-long-id-for-location

# Set the URL for the Backstage, the authentication token, and Skip-TLS settings
$ %s get locations --backstage-url=https://my-rhdh.com --backstage-token=my-token --backstage-skip-tls=true
`

	getComponentsExample = `
# Retrieve the Backstage Catalog for resources related to AI Models, where being AI related is determined by the 
# 'type' being set to 'model-server'
$ %s get components [args...]

# Set the URL for the Backstage, the authentication token, and Skip-TLS settings
$ %s get components --backstage-url=https://my-rhdh.com --backstage-token=my-token --backstage-skip-tls=true

# Retrieve a specific set of AI related Components by namespace:name
$ %s get components default:my-component default:your-component

# Retrieve a set of AI Components where the provided list of tags match (order of tags disregarded)
$ %s get components genai vllm --use-params-as-tags=true

# Retrieve a set of Components which have any of the provided list of tags
$ %s get components gen-ai --use-params-as-tags=true --use-any-subset=true
`

	getResourcesExample = `
# Retrieve the Backstage Catalog for resources related to AI Models, where being AI related is determined by the 
# 'type' being set to 'ai-model'
$ %s get resources [args...]

# Set the URL for the Backstage, the authentication token, and Skip-TLS settings
$ %s get resources --backstage-url=https://my-rhdh.com --backstage-token=my-token --backstage-skip-tls=true

# Retrieve a specific set of AI related Resources by namespace:name
$ %s get resources default:my-component default:your-component

# Retrieve a set of AI Resources where the provided list of tags match (order of tags disregarded)
$ %s get resources genai vllm --use-params-as-tags=true

# Retrieve a set of AI Resources which have any of the provided list of tags
$ %s get resources gen-ai --use-params-as-tags=true --use-any-subset=true
`

	getApisExample = `
# Retrieve the Backstage Catalog for APIs related to AI Models, where being AI related is determined by the 
# 'type' being set to 'model-service-api'
$ %s get apis [args...]

# Set the URL for the Backstage, the authentication token, and Skip-TLS settings
$ %s get locations --backstage-url=https://my-rhdh.com --backstage-token=my-token --backstage-skip-tls=true

# Retrieve a specific set of AI related APIs by namespace:name
$ %s get apis default:my-component default:your-component

# Retrieve a set of AI APIs where the provided list of tags match (order of tags disregarded)
$ %s get apis genai vllm --use-params-as-tags=true

# Retrieve a set of AI APIs which have any of the provided list of tags
$ %s get apis gen-ai --use-params-as-tags=true --use-any-subset=true
`
)

// NewCmd create a new root command, linking together all sub-commands organized by groups.
func NewCmd() *cobra.Command {
	cfg := &config.Config{}
	bkstgAI := &cobra.Command{
		Use:     util.ApplicationName,
		Long:    "Backstage AI is a command line tool that facilitates management of AI related Entities in the Backstage Catalog.",
		Example: strings.ReplaceAll(bkstgAIExample, "%s", util.ApplicationName),
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
		Example: strings.ReplaceAll(newModelExample, "%s", util.ApplicationName),
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
		Example: strings.ReplaceAll(getExample, "%s", util.ApplicationName),
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	deleteModel := &cobra.Command{
		Use:     "delete-model",
		Long:    "delete-model removes the Backstage Catalog for Entities corresponding to the provided location ID",
		Aliases: []string{"delete", "dm", "del", "d", "delete-models"},
		Example: strings.ReplaceAll(deleteModelExample, "%s", util.ApplicationName),
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
		Example: strings.ReplaceAll(importModelExample, "%s", util.ApplicationName),
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
		Example: strings.ReplaceAll(getEntitiesExample, "%s", util.ApplicationName),
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
		Example: strings.ReplaceAll(getLocationsExample, "%s", util.ApplicationName),
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
		Example: strings.ReplaceAll(getComponentsExample, "%s", util.ApplicationName),
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
		Example: strings.ReplaceAll(getResourcesExample, "%s", util.ApplicationName),
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
		Example: strings.ReplaceAll(getApisExample, "%s", util.ApplicationName),
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
