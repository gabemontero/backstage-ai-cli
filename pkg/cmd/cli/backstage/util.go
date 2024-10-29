package backstage

import (
	"net/url"
	"strings"
)

func buildKeys(args ...string) map[string][]string {
	keys := map[string][]string{}
	for _, arg := range args {
		array := strings.Split(arg, ":")
		if len(array) == 1 {
			arr := keys[DEFAULT_NS]
			arr = append(arr, arg)
			keys[DEFAULT_NS] = arr
			continue
		}
		arr := keys[array[0]]
		arr = append(arr, array[1])
		keys[array[0]] = arr
	}
	return keys
}

func (b *BackstageRESTClientWrapper) pullSavedArgsFromQueryParams(qparms *url.Values) []string {
	var argsArr []string
	if b.Tags && !b.Subset {
		argsStr := qparms.Get("metadata.tags")
		argsArr = strings.Split(argsStr, " ")
		qparms.Del("metadata.tags")
	}
	return argsArr
}

func tagsMatch(args, tags []string) bool {
	for i, tag := range tags {
		if args[i] != tag {
			return false
		}
	}
	return true
}
