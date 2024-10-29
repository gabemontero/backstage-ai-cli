package kserve

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	clibkstg "github.com/gabemontero/backstage-ai-cli/pkg/cmd/cli/backstage"
	"github.com/gabemontero/backstage-ai-cli/pkg/config"
	"github.com/gabemontero/backstage-ai-cli/pkg/util"
	serverapiv1beta1 "github.com/kserve/kserve/pkg/apis/serving/v1beta1"
	servingv1beta1 "github.com/kserve/kserve/pkg/client/clientset/versioned/typed/serving/v1beta1"
	"github.com/spf13/cobra"
	"github.com/tdabasinskas/go-backstage/v2/backstage"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"strings"
)

type commonPopulator struct {
	owner     string
	lifecycle string
	is        *serverapiv1beta1.InferenceService
}

func (pop *commonPopulator) GetOwner() string {
	return pop.owner
}

func (pop *commonPopulator) GetLifecycle() string {
	return pop.lifecycle
}

func (pop *commonPopulator) GetName() string {
	return fmt.Sprintf("%s_%s", pop.is.Namespace, pop.is.Name)
}
func (pop *commonPopulator) GetDescription() string {
	return fmt.Sprintf("KServe instance %s:%s", pop.is.Namespace, pop.is.Name)
}
func (pop *commonPopulator) GetDisplayName() string {
	return fmt.Sprintf("KServe instance %s:%s", pop.is.Namespace, pop.is.Name)
}

func (pop *commonPopulator) GetLinks() []backstage.EntityLink {
	links := []backstage.EntityLink{}
	if pop.is.Status.URL != nil {
		links = append(links, backstage.EntityLink{
			URL:   pop.is.Status.URL.String(),
			Title: clibkstg.LINK_API_URL,
			Type:  clibkstg.LINK_TYPE_WEBSITE,
			Icon:  clibkstg.LINK_ICON_WEBASSET,
		})
	}
	for componentType, componentStatus := range pop.is.Status.Components {

		if componentStatus.URL != nil {
			links = append(links, backstage.EntityLink{
				URL:   componentStatus.URL.String() + "/docs",
				Title: string(componentType) + " FastAPI URL",
				Icon:  clibkstg.LINK_ICON_WEBASSET,
				Type:  clibkstg.LINK_TYPE_WEBSITE,
			})
			links = append(links, backstage.EntityLink{
				URL:   componentStatus.URL.String(),
				Title: string(componentType) + " model serving URL",
				Icon:  clibkstg.LINK_ICON_WEBASSET,
				Type:  clibkstg.LINK_TYPE_WEBSITE,
			})
		}
		if componentStatus.RestURL != nil {
			links = append(links, backstage.EntityLink{
				URL:   componentStatus.RestURL.String(),
				Title: string(componentType) + " REST model serving URL",
				Icon:  clibkstg.LINK_ICON_WEBASSET,
				Type:  clibkstg.LINK_TYPE_WEBSITE,
			})
		}
		if componentStatus.GrpcURL != nil {
			links = append(links, backstage.EntityLink{
				URL:   componentStatus.GrpcURL.String(),
				Title: string(componentType) + " GRPC model serving URL",
				Icon:  clibkstg.LINK_ICON_WEBASSET,
				Type:  clibkstg.LINK_TYPE_WEBSITE,
			})
		}
	}
	return links
}

func (pop *commonPopulator) GetTags() []string {
	tags := []string{}
	predictor := pop.is.Spec.Predictor
	tag := ""
	// one and only one predictor spec can be set
	switch {
	case predictor.SKLearn != nil:
		tag = "sklearn"
	case predictor.XGBoost != nil:
		tag = "xgboost"
	case predictor.Tensorflow != nil:
		tag = "tensorflow"
	case predictor.PyTorch != nil:
		tag = "pytorch"
	case predictor.Triton != nil:
		tag = "triton"
	case predictor.ONNX != nil:
		tag = "onnx"
	case predictor.HuggingFace != nil:
		tag = "huggingface"
	case predictor.PMML != nil:
		tag = "pmml"
	case predictor.LightGBM != nil:
		tag = "lightgbm"
	case predictor.Paddle != nil:
		tag = "paddle"
	case predictor.Model != nil:
		modelFormat := predictor.Model.ModelFormat
		tag = modelFormat.Name
		if modelFormat.Version != nil {
			tag = tag + "-" + *modelFormat.Version
		}
		tag = strings.ToLower(tag)
	}
	tags = append(tags, tag)
	explainer := pop.is.Spec.Explainer
	if explainer != nil && explainer.ART != nil {
		tags = append(tags, strings.ToLower(string(explainer.ART.Type)))
	}
	return tags
}

func (pop *commonPopulator) GetProvidedAPIs() []string {
	return []string{fmt.Sprintf("%s_%s", pop.is.Namespace, pop.is.Name)}
}

type componentPopulator struct {
	commonPopulator
}

func (pop *componentPopulator) GetDependsOn() []string {
	return []string{fmt.Sprintf("resource:%s_%s", pop.is.Namespace, pop.is.Name), fmt.Sprintf("api:%s_%s", pop.is.Namespace, pop.is.Name)}
}

func (pop *componentPopulator) GetTechdocRef() string {
	return "./"
}

type resourcePopulator struct {
	commonPopulator
}

func (pop *resourcePopulator) GetDependencyOf() []string {
	return []string{fmt.Sprintf("component:%s_%s", pop.is.Namespace, pop.is.Name)}
}

func (pop *resourcePopulator) GetTechdocRef() string {
	return "resource/"
}

type apiPopulator struct {
	commonPopulator
}

func (pop *apiPopulator) GetDependencyOf() []string {
	return []string{fmt.Sprintf("component:%s_%s", pop.is.Namespace, pop.is.Name)}
}

func (pop *apiPopulator) GetDefinition() string {
	if pop.is.Status.URL == nil {
		return ""
	}
	defBytes, _ := util.FetchURL(pop.is.Status.URL.String() + "/openapi.json")
	dst := bytes.Buffer{}
	json.Indent(&dst, defBytes, "", "    ")
	return dst.String()
}

func (pop *apiPopulator) GetTechdocRef() string {
	return "api/"
}

func NewCmd(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kserve",
		Short: "KServe related API",
		Long:  "Interact with KServe related instances on a K8s cluster to manage AI related catalog entities in a Backstage instance.",
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

			var servingClient servingv1beta1.ServingV1beta1Interface
			if kubeconfig, err := util.GetK8sConfig(cfg); err != nil {
				klog.Errorf("ERROR: problem with kubeconfig: %v\n", err)
				klog.Flush()
				return
			} else {
				servingClient = util.GetKServeClient(kubeconfig)
			}

			namespace := cfg.Namespace
			if len(namespace) == 0 {
				namespace = util.GetCurrentProject()
			}
			if len(ids) != 0 {
				for _, id := range ids {
					is, err := servingClient.InferenceServices(namespace).Get(context.Background(), id, metav1.GetOptions{})
					if err != nil {
						klog.Errorf("ERROR: inference service retrieval error for %s:%s: %s\n", namespace, id, err.Error())
						klog.Flush()
						return
					}

					err = callBackstagePrinters(owner, lifecycle, is)
					if err != nil {
						return
					}
				}
			} else {
				isl, err := servingClient.InferenceServices(namespace).List(context.Background(), metav1.ListOptions{})
				if err != nil {
					klog.Errorf("ERROR: inference service retrieval error for %s: %s\n", namespace, err.Error())
					klog.Flush()
					return
				}
				for _, is := range isl.Items {
					err = callBackstagePrinters(owner, lifecycle, &is)
					if err != nil {
						return
					}
				}
			}

		},
	}

	return cmd
}

func callBackstagePrinters(owner, lifecycle string, is *serverapiv1beta1.InferenceService) error {
	compPop := componentPopulator{}
	compPop.owner = owner
	compPop.lifecycle = lifecycle
	compPop.is = is
	err := clibkstg.PrintComponent(&compPop)
	if err != nil {
		return err
	}

	resPop := resourcePopulator{}
	resPop.owner = owner
	resPop.lifecycle = lifecycle
	resPop.is = is
	err = clibkstg.PrintResource(&resPop)
	if err != nil {
		return err
	}

	apiPop := apiPopulator{}
	apiPop.owner = owner
	apiPop.lifecycle = lifecycle
	apiPop.is = is
	err = clibkstg.PrintAPI(&apiPop)
	return err
}
