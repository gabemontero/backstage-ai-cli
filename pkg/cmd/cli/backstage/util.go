package backstage

import (
	"net/url"
	"sort"
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
	if len(args) != len(tags) {
		return false
	}
	// we don't require exact order with the set of tags specified so we sort the two arrays to facilitate the compare
	sort.Strings(args)
	sort.Strings(tags)
	for i, tag := range tags {
		if args[i] != tag {
			return false
		}
	}
	return true
}

func updateQParams(subset bool, filterValue string, args []string, qparams *url.Values) *url.Values {
	if subset {
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
	return qparams
}
