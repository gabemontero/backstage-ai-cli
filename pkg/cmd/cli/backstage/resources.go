package backstage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tdabasinskas/go-backstage/v2/backstage"
	"net/url"
	"strings"
)

type listResources struct {
	Items      []backstage.ResourceEntityV1alpha1 `json:"items" yaml:"items"`
	TotalItems int                                `json:"totalItems" yaml:"totalItems"`
	PageInfo   interface{}                        `json:"pageInfo" yaml:"pageInfo"`
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
		filteredResources := []backstage.ResourceEntityV1alpha1{}
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
		if b.Subset {
			for _, arg := range args {
				filterValue = filterValue + ",metadata.tags=" + arg
			}
			qparams.Set("filter", filterValue)
		} else {
			//TODO could not determine single query parameter format that resulted in returning
			// a list of entities whose `metadata.tags` array directly matched a provided list of args;
			// for now, we capture the arg list in the query params
			qparams.Add("metadata.tags", strings.Join(args, " "))
		}

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
