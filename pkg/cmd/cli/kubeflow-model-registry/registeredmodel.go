package kubeflow_model_registry

import (
	"fmt"
	"github.com/kubeflow/model-registry/pkg/openapi"
)

// TODO need params
func (k *KubeFlowRESTClientWrapper) CreateRegisteredModel(registeredModelID string) (string, error) {
	modelRegistry := openapi.RegisteredModel{}

	modelRegistry = openapi.RegisteredModel{}
	modelRegistry.Name = registeredModelID
	desc := "Description for " + registeredModelID
	modelRegistry.Description = &desc
	return k.postToModelRegistry(k.RootURL+CREATE_REG_MODEL_URI, marshalBody(modelRegistry))
}

// TODO need params
func (k *KubeFlowRESTClientWrapper) PatchRegisteredModel(registeredModelID string) (string, error) {
	modelRegistry := openapi.RegisteredModel{}

	modelRegistry = openapi.RegisteredModel{}
	modelRegistry.Name = registeredModelID
	desc := "Description for " + registeredModelID
	modelRegistry.Description = &desc
	return k.patchToModelRegistry(k.RootURL+fmt.Sprintf(PATCH_REG_MODEL_URI, registeredModelID), marshalBody(modelRegistry))
}

func (k *KubeFlowRESTClientWrapper) ListRegisteredModels() (string, error) {
	return k.getFromModelRegistry(k.RootURL + LIST_REG_MODEL_URI)
}

func (k *KubeFlowRESTClientWrapper) GetRegisteredModel(registeredModelID string) (string, error) {
	return k.getFromModelRegistry(k.RootURL + fmt.Sprintf(GET_REG_MODEL_URI, registeredModelID))
}
