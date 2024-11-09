package backstage

import (
	"github.com/gabemontero/backstage-ai-cli/pkg/util"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

type CommonPopulator interface {
	GetOwner() string
	GetLifecycle() string
	GetName() string
	GetDescription() string
	GetLinks() []EntityLink
	GetTags() []string
	GetDisplayName() string
	GetProvidedAPIs() []string
	GetTechdocRef() string
}

type ComponentPopulator interface {
	CommonPopulator
	GetDependsOn() []string
}

type ResourcePopulator interface {
	CommonPopulator
	GetDependencyOf() []string
}

type APIPopulator interface {
	CommonPopulator
	GetDefinition() string
	GetDependencyOf() []string
}

func PrintComponent(pop ComponentPopulator, cmd *cobra.Command) error {
	component := &ComponentEntityV1alpha1{
		Kind:       "Component",
		ApiVersion: VERSION,
		Entity:     buildEntity("Component", pop),
	}
	component.Entity.Metadata.Annotations = map[string]string{TECHDOC_REFS: pop.GetTechdocRef()}
	component.Metadata = component.Entity.Metadata
	component.Spec = &ComponentEntityV1alpha1Spec{
		Type:         COMPONENT_TYPE,
		Lifecycle:    pop.GetLifecycle(),
		Owner:        "user:" + pop.GetOwner(),
		ProvidesApis: pop.GetProvidedAPIs(),
		DependsOn:    pop.GetDependsOn(),
		//TODO this version of the converted to Golang does not have `profile` with `displayName`
	}
	err := util.PrintYaml(component, true, cmd)
	if err != nil {
		klog.Errorf("ERROR: converting component to yaml and printing: %s, %#v", err.Error(), component)
		return err
	}
	return nil
}

func PrintResource(pop ResourcePopulator, cmd *cobra.Command) error {
	resource := &ResourceEntityV1alpha1{
		Kind:       "Resource",
		ApiVersion: VERSION,
		Entity:     buildEntity("Resource", pop),
	}
	resource.Entity.Metadata.Annotations = map[string]string{TECHDOC_REFS: pop.GetTechdocRef()}
	resource.Metadata = resource.Entity.Metadata
	resource.Spec = &ResourceEntityV1alpha1Spec{
		Type:         RESOURCE_TYPE,
		Owner:        "user:" + pop.GetOwner(),
		Lifecycle:    pop.GetLifecycle(),
		ProvidesApis: pop.GetProvidedAPIs(),
		DependencyOf: pop.GetDependencyOf(),
	}
	err := util.PrintYaml(resource, true, cmd)
	if err != nil {
		klog.Errorf("ERROR: converting resource to yaml and printing: %s, %#v", err.Error(), resource)
		return err
	}
	return nil
}

func PrintAPI(pop APIPopulator, cmd *cobra.Command) error {
	api := &ApiEntityV1alpha1{
		Kind:       "API",
		ApiVersion: VERSION,
		Entity:     buildEntity("API", pop),
	}
	api.Entity.Metadata.Annotations = map[string]string{TECHDOC_REFS: pop.GetTechdocRef()}
	api.Metadata = api.Entity.Metadata
	api.Spec = &ApiEntityV1alpha1Spec{
		Type:         API_TYPE,
		Lifecycle:    pop.GetLifecycle(),
		Owner:        "user:" + pop.GetOwner(),
		Definition:   pop.GetDefinition(),
		DependencyOf: pop.GetDependencyOf(),
	}
	err := util.PrintYaml(api, false, cmd)
	if err != nil {
		klog.Errorf("ERROR: converting api to yaml and printing: %s, %#v", err.Error(), api)
		return err
	}
	return nil
}

func buildEntity(kind string, pop CommonPopulator) Entity {
	entity := Entity{
		Kind:       kind,
		ApiVersion: VERSION,
		Metadata: EntityMeta{
			Name:        pop.GetName(),
			Description: pop.GetDescription(),
			Tags:        pop.GetTags(),
			Links:       pop.GetLinks(),
		},
	}
	return entity
}
