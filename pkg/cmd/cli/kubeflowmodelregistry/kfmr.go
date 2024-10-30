package kubeflowmodelregistry

import (
	"fmt"
	clibkstg "github.com/gabemontero/backstage-ai-cli/pkg/cmd/cli/backstage"
	"github.com/gabemontero/backstage-ai-cli/pkg/config"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"github.com/spf13/cobra"
	"github.com/tdabasinskas/go-backstage/v2/backstage"
	"k8s.io/klog/v2"
)

func NewCmd(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "kubeflow",
		Aliases: []string{"kf"},
		Short:   "Kubeflow Model Registry related API",
		Long:    "Interact with the Kubeflow Model Registry REST API as part of managing AI related catalog entities in a Backstage instance.",
		Run: func(cmd *cobra.Command, args []string) {
			ids := []string{}

			if len(args) < 2 {
				klog.Errorf("ERROR: need to specify an owner and lifecycle setting")
				klog.Flush()
				return
			}
			owner := args[0]
			lifecycle := args[1]

			if len(args) > 2 {
				ids = args[2:]
			}

			kfmr := SetupKubeflowRESTClient(cfg)

			if len(ids) == 0 {
				var err error
				var rms []openapi.RegisteredModel
				rms, err = kfmr.ListRegisteredModels()
				if err != nil {
					klog.Errorf("ERROR: list registered models error: %s", err.Error())
					klog.Flush()
					return
				}
				for _, rm := range rms {
					var mvs []openapi.ModelVersion
					var mas map[string][]openapi.ModelArtifact
					mvs, mas, err = callKubeflowREST(*rm.Id, kfmr)
					if err != nil {
						return
					}
					err = callBackstagePrinters(owner, lifecycle, &rm, mvs, mas)
					if err != nil {
						klog.Errorf("ERROR: print model catalog: %s", err.Error())
						klog.Flush()
						return
					}
				}
			} else {
				for _, id := range ids {
					rm, err := kfmr.GetRegisteredModel(id)
					if err != nil {
						klog.Errorf("ERROR: get registered model error for %s: %s", id, err.Error())
						klog.Flush()
						return
					}
					var mvs []openapi.ModelVersion
					var mas map[string][]openapi.ModelArtifact
					mvs, mas, err = callKubeflowREST(*rm.Id, kfmr)
					if err != nil {
						klog.Errorf("ERROR: get model version/artifact error for %s: %s", id, err.Error())
						klog.Flush()
						return
					}
					err = callBackstagePrinters(owner, lifecycle, rm, mvs, mas)
				}
			}

		},
	}

	return cmd
}

func callKubeflowREST(id string, kfmr *KubeFlowRESTClientWrapper) (mvs []openapi.ModelVersion, ma map[string][]openapi.ModelArtifact, err error) {
	mvs, err = kfmr.ListModelVersions(id)
	if err != nil {
		klog.Errorf("ERROR: error list model versions for %s: %s", id, err.Error())
		return
	}
	ma = map[string][]openapi.ModelArtifact{}
	for _, mv := range mvs {
		var v []openapi.ModelArtifact
		v, err = kfmr.ListModelArtifacts(*mv.Id)
		if err != nil {
			klog.Errorf("ERROR error list model artifacts for %s:%s: %s", id, *mv.Id, err.Error())
			return
		}
		if len(v) == 0 {
			v, err = kfmr.ListModelArtifacts(id)
			if err != nil {
				klog.Errorf("ERROR error list model artifacts for %s:%s: %s", id, *mv.Id, err.Error())
				return
			}
		}
		ma[*mv.Id] = v
	}
	return
}

func callBackstagePrinters(owner, lifecycle string, rm *openapi.RegisteredModel, mv []openapi.ModelVersion, ma map[string][]openapi.ModelArtifact) error {
	compPop := componentPopulator{}
	compPop.owner = owner
	compPop.lifecycle = lifecycle
	compPop.registeredModel = rm
	compPop.modelVersions = mv
	compPop.modelArtifacts = ma
	err := clibkstg.PrintComponent(&compPop)
	if err != nil {
		return err
	}

	resPop := resourcePopulator{}
	resPop.owner = owner
	resPop.lifecycle = lifecycle
	resPop.registeredModel = rm
	resPop.modelVersions = mv
	resPop.modelArtifacts = ma
	err = clibkstg.PrintResource(&resPop)
	if err != nil {
		return err
	}

	apiPop := apiPopulator{}
	apiPop.owner = owner
	apiPop.lifecycle = lifecycle
	apiPop.registeredModel = rm
	apiPop.modelVersions = mv
	apiPop.modelArtifacts = ma
	err = clibkstg.PrintAPI(&apiPop)
	return err
}

type commonPopulator struct {
	owner           string
	lifecycle       string
	registeredModel *openapi.RegisteredModel
	modelVersions   []openapi.ModelVersion
	modelArtifacts  map[string][]openapi.ModelArtifact
}

func (pop *commonPopulator) GetOwner() string {
	if pop.registeredModel.Owner != nil {
		return *pop.registeredModel.Owner
	}
	return pop.owner
}

func (pop *commonPopulator) GetLifecycle() string {
	return pop.lifecycle
}

func (pop *commonPopulator) GetName() string {
	return pop.registeredModel.Name
}

func (pop *commonPopulator) GetDescription() string {
	if pop.registeredModel.Description != nil {
		return *pop.registeredModel.Description
	}
	return ""
}

func (pop *commonPopulator) GetDisplayName() string {
	if pop.registeredModel.ExternalId != nil {
		return *pop.registeredModel.ExternalId
	}
	return fmt.Sprintf("%s model from KubeFlow Model Registry", pop.registeredModel.Name)
}

// TODO won't have API until KubeFlow Model Registry gets us inferenceservice endpoints
func (pop *commonPopulator) GetProvidedAPIs() []string {
	return []string{}
}

type componentPopulator struct {
	commonPopulator
}

// TODO Until we get the inferenceservice endpoint URL associated with the model registry related API won't have component links
func (pop *componentPopulator) GetLinks() []backstage.EntityLink {
	return []backstage.EntityLink{}
}

func (pop *componentPopulator) GetTags() []string {
	tags := []string{}
	for key, value := range pop.registeredModel.GetCustomProperties() {
		tags = append(tags, fmt.Sprintf("%s:%v", key, value.GetActualInstance()))
	}

	return tags
}

func (pop *componentPopulator) GetDependsOn() []string {
	depends := []string{}
	for _, mv := range pop.modelVersions {
		depends = append(depends, "resource:"+mv.Name)
	}
	for _, mas := range pop.modelArtifacts {
		for _, ma := range mas {
			depends = append(depends, "api:"+*ma.Name)
		}
	}
	return depends
}

func (pop *componentPopulator) GetTechdocRef() string {
	return "./"
}

type resourcePopulator struct {
	commonPopulator
}

func (pop *resourcePopulator) GetTechdocRef() string {
	return "resource/"
}

func (pop *resourcePopulator) GetLinks() []backstage.EntityLink {
	links := []backstage.EntityLink{}
	for _, value := range pop.modelArtifacts {
		for _, ma := range value {
			if ma.Uri != nil {
				links = append(links, backstage.EntityLink{
					URL:   *ma.Uri,
					Title: ma.GetDescription(),
					Icon:  clibkstg.LINK_ICON_WEBASSET,
					Type:  clibkstg.LINK_TYPE_WEBSITE,
				})
			}
		}
	}
	return links
}

func (pop *resourcePopulator) GetTags() []string {
	tags := []string{}
	for _, mv := range pop.modelVersions {
		for key := range mv.GetCustomProperties() {
			tags = append(tags, key)
		}
	}
	for _, value := range pop.modelArtifacts {
		for _, ma := range value {
			for k := range ma.GetCustomProperties() {
				tags = append(tags, k)
			}
		}
	}
	return tags
}

func (pop *resourcePopulator) GetDependencyOf() []string {
	return []string{fmt.Sprintf("component:%s", pop.registeredModel.Name)}
}

// TODO Until we get the inferenceservice endpoint URL associated with the model registry related API won't have much for Backstage API here
type apiPopulator struct {
	commonPopulator
}

func (pop *apiPopulator) GetDependencyOf() []string {
	return []string{fmt.Sprintf("component:%s", pop.registeredModel.Name)}
}

func (pop *apiPopulator) GetDefinition() string {
	// definition must be set to something to pass backstage validation
	return "no-definition-yet"
}

func (pop *apiPopulator) GetTechdocRef() string {
	return "api/"
}

func (pop *apiPopulator) GetTags() []string {
	return []string{}
}

func (pop *apiPopulator) GetLinks() []backstage.EntityLink {
	return []backstage.EntityLink{}
}
