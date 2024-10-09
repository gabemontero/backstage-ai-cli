package config

type Config struct {
	// K8S related
	Kubeconfig string
	Namespace  string

	// Kubeflow Model Registry Related
	KubeFlowSkipTLS bool
	KubeFlowToken   string
	KubeFlowURL     string
}
