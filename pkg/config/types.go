package config

type Config struct {
	// K8S related
	Kubeconfig string
	Namespace  string

	// Cross "store" related
	StoreURL     string
	StoreToken   string
	StoreSkipTLS bool

	// Backstage related
	BackstageSkipTLS bool
	BackstageToken   string
	BackstageURL     string

	// new-model related
	DeleteAll              bool
	ConfigMapNS            string
	ConfigMapName          string
	Owner                  string
	Lifecycle              string
	ComponentTags          []string
	ResourceTags           map[string][]string
	APITags                []string
	ComponentLinks         map[string]Link
	ResourceLinks          map[string]map[string]Link
	APILinks               map[string]Link
	ComponentTechDockRef   string
	ResourceTechDockRef    map[string]string
	APITechDockRef         string
	MultiEntryOutputPrefix string

	// fetch-model related
	ParamsAsTags   bool
	AnySubsetWorks bool
}

type Link struct {
	Title string
	Type  string
	Icon  string
}
