package backstage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
)

type listComponents struct {
	Items      []ComponentEntityV1alpha1 `json:"items" yaml:"items"`
	TotalItems int                       `json:"totalItems" yaml:"totalItems"`
	PageInfo   interface{}               `json:"pageInfo" yaml:"pageInfo"`
}

func (b *BackstageRESTClientWrapper) ListComponents(qparms *url.Values) (string, error) {
	//TODO remove this post query filter logic if an exact query parameter check for the 'metadata.tags' array is determined
	argsArr := b.pullSavedArgsFromQueryParams(qparms)

	str, err := b.getWithKindParamFromBackstage(b.RootURL+QUERY_URI, qparms)
	if err != nil {
		return str, err
	}

	buf := []byte(str)

	lc := &listComponents{}
	err = json.Unmarshal(buf, lc)
	if err != nil {
		return str, err
	}

	//TODO remove this post query filter logic if an exact query parameter check for the 'metadata.tags' array is determined
	if b.Tags && !b.Subset {
		filteredComponents := []ComponentEntityV1alpha1{}
		for _, component := range lc.Items {
			if tagsMatch(argsArr, component.Metadata.Tags) {
				filteredComponents = append(filteredComponents, component)
			}
		}
		lc.Items = filteredComponents
	}

	buf, err = json.MarshalIndent(lc.Items, "", "    ")
	return string(buf), err
}

func (b *BackstageRESTClientWrapper) GetComponent(args ...string) (string, error) {
	if len(args) == 0 {
		return b.ListComponents(&url.Values{
			"filter": []string{"kind=component", "spec.type=model-server"},
		})
	}
	keys := buildKeys(args...)
	buffer := &bytes.Buffer{}
	for namespace, names := range keys {
		for _, name := range names {
			str, err := b.getFromBackstage(b.RootURL + fmt.Sprintf(COMPONENT_URI, namespace, name))
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
