package backstage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
)

type listResources struct {
	Items      []ResourceEntityV1alpha1 `json:"items" yaml:"items"`
	TotalItems int                      `json:"totalItems" yaml:"totalItems"`
	PageInfo   interface{}              `json:"pageInfo" yaml:"pageInfo"`
}

func (b *BackstageRESTClientWrapper) ListResources(qparms *url.Values) (string, error) {
	//TODO remove this post query filter logic if an exact query parameter check for the 'metadata.tags' array is determined
	argsArr := b.pullSavedArgsFromQueryParams(qparms)

	str, err := b.getWithKindParamFromBackstage(b.RootURL+QUERY_URI, qparms)
	if err != nil {
		return "", err
	}

	buf := []byte(str)

	lr := &listResources{}
	err = json.Unmarshal(buf, lr)
	if err != nil {
		return str, err
	}

	//TODO remove this post query filter logic if an exact query parameter check for the 'metadata.tags' array is determined
	if b.Tags && !b.Subset {
		filteredResources := []ResourceEntityV1alpha1{}
		for _, resource := range lr.Items {
			if tagsMatch(argsArr, resource.Metadata.Tags) {
				filteredResources = append(filteredResources, resource)
			}
		}
		lr.Items = filteredResources
	}

	buf, err = json.MarshalIndent(lr.Items, "", "    ")
	return string(buf), err
}

func (b *BackstageRESTClientWrapper) GetResource(args ...string) (string, error) {
	// example 'filter' value from swagger doc:  'kind=component,metadata.annotations.backstage.io/orphan=true'
	filterValue := "kind=resource,spec.type=ai-model"
	qparams := &url.Values{
		"filter": []string{filterValue},
	}
	if len(args) == 0 {
		return b.ListResources(qparams)
	}

	if b.Tags {
		qparams = updateQParams(b.Subset, filterValue, args, qparams)
		return b.ListResources(qparams)
	}

	keys := buildKeys(args...)
	buffer := &bytes.Buffer{}
	for namespace, names := range keys {
		for _, name := range names {
			str, err := b.getFromBackstage(b.RootURL + fmt.Sprintf(RESOURCE_URI, namespace, name))
			if err != nil {
				return buffer.String(), err
			}
			buf := []byte(str)
			err = json.Indent(buffer, buf, "", "    ")
			buffer.WriteString("\n")
		}
	}

	return buffer.String(), nil
}
