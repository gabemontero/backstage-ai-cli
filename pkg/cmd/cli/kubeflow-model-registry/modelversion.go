package kubeflow_model_registry

import (
	"fmt"
	"github.com/kubeflow/model-registry/pkg/openapi"
)

// TODO need params
func (k *KubeFlowRESTClientWrapper) CreateModelVersion(registeredModelID, modelVersionID string) (string, error) {
	modelVersion := openapi.NewModelVersion(modelVersionID, registeredModelID)

	desc := "Description for " + modelVersionID
	modelVersion.Description = &desc
	return k.postToModelRegistry(k.RootURL+fmt.Sprintf(CREATE_MODEL_VERSION_URI, registeredModelID), marshalBody(modelVersion))
}

// TODO need params
func (k *KubeFlowRESTClientWrapper) PatchModelVersion(registeredModelID, modelVersionID string) (string, error) {
	modelVersion := openapi.NewModelVersion(modelVersionID, registeredModelID)

	desc := "Description for " + modelVersionID
	modelVersion.Description = &desc
	return k.patchToModelRegistry(k.RootURL+fmt.Sprintf(PATCH_MODEL_VERSION_URI, modelVersionID), marshalBody(modelVersion))
}

func (k *KubeFlowRESTClientWrapper) ListModelVersions() (string, error) {
	return k.getFromModelRegistry(k.RootURL + LIST_MODEL_VERSION_URI)
}

func (k *KubeFlowRESTClientWrapper) GetModelVersion(modelVersionID string) (string, error) {
	return k.getFromModelRegistry(k.RootURL + fmt.Sprintf(GET_MODEL_VERSION_URI, modelVersionID))
}
