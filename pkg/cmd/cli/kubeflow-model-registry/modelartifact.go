package kubeflow_model_registry

import (
	"fmt"
	"github.com/kubeflow/model-registry/pkg/openapi"
)

// TODO need params
func (k *KubeFlowRESTClientWrapper) CreateModelArtifact(modelVersionID, modelArtifactID string) (string, error) {
	modelArtifact := &openapi.ModelArtifact{}
	modelArtifact.Name = &modelArtifactID
	desc := "Description for " + modelArtifactID
	modelArtifact.Description = &desc
	modelArtifact.ArtifactType = "model-artifact"
	return k.postToModelRegistry(k.RootURL+fmt.Sprintf(CREATE_MODEL_ART_URI, modelVersionID), marshalBody(modelArtifact))
}

// TODO need params
func (k *KubeFlowRESTClientWrapper) PatchModelArtifact(modelArtifactID string) (string, error) {
	modelArtifact := &openapi.ModelArtifact{}
	modelArtifact.Name = &modelArtifactID
	desc := "Description for " + modelArtifactID
	modelArtifact.Description = &desc
	modelArtifact.ArtifactType = "model-artifact"
	return k.patchToModelRegistry(k.RootURL+fmt.Sprintf(PATCH_MODEL_ART_URI, modelArtifactID), marshalBody(modelArtifact))
}

func (k *KubeFlowRESTClientWrapper) ListModelArtifacts() (string, error) {
	return k.getFromModelRegistry(k.RootURL + LIST_MODEL_ART_URI)
}

func (k *KubeFlowRESTClientWrapper) GetModelArtifact(modelVersionID string) (string, error) {
	return k.getFromModelRegistry(k.RootURL + fmt.Sprintf(GET_MODEL_ART_URI, modelVersionID))
}
