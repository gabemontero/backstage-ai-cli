package backstage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
)

type listAPIs struct {
	Items      []ApiEntityV1alpha1 `json:"items" yaml:"items"`
	TotalItems int                 `json:"totalItems" yaml:"totalItems"`
	PageInfo   interface{}         `json:"pageInfo" yaml:"pageInfo"`
}

func (b *BackstageRESTClientWrapper) ListAPIs(qparms *url.Values) (string, error) {
	//TODO remove this post query filter logic if an exact query parameter check for the 'metadata.tags' array is determined
	argsArr := b.pullSavedArgsFromQueryParams(qparms)

	str, err := b.getWithKindParamFromBackstage(b.RootURL+QUERY_URI, qparms)
	if err != nil {
		return "", err
	}

	buf := []byte(str)

	la := &listAPIs{}
	err = json.Unmarshal(buf, la)
	if err != nil {
		return str, err
	}

	//TODO remove this post query filter logic if an exact query parameter check for the 'metadata.tags' array is determined
	if b.Tags && !b.Subset {
		filteredAPIs := []ApiEntityV1alpha1{}
		for _, api := range la.Items {
			if tagsMatch(argsArr, api.Metadata.Tags) {
				filteredAPIs = append(filteredAPIs, api)
			}
		}
		la.Items = filteredAPIs
	}

	buf, err = json.MarshalIndent(la.Items, "", "    ")
	return string(buf), err
}

func (b *BackstageRESTClientWrapper) GetAPI(args ...string) (string, error) {
	// example 'filter' value from swagger doc:  'kind=component,metadata.annotations.backstage.io/orphan=true'
	filterValue := "kind=api,spec.type=openapi"
	qparams := &url.Values{
		"filter": []string{filterValue},
	}
	if len(args) == 0 {
		return b.ListAPIs(qparams)
	}

	if b.Tags {
		qparams = updateQParams(b.Subset, filterValue, args, qparams)
		return b.ListAPIs(qparams)
	}

	keys := buildKeys(args...)
	buffer := &bytes.Buffer{}
	for namespace, names := range keys {
		for _, name := range names {
			str, err := b.getFromBackstage(b.RootURL + fmt.Sprintf(API_URI, namespace, name))
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
