package elmohttp

import (
	"net/url"
	"sort"
	"strings"

	elmo "github.com/okke/elmo/core"
)

func addParametersToPath(path string, parameters elmo.DictionaryValue) string {

	if len(parameters.Keys()) == 0 {
		return path
	}

	var sb strings.Builder
	sb.WriteString(path)
	if len(path) == 0 || path[len(path)-1] != '?' {
		sb.WriteRune('?')
	}

	keys := parameters.Keys()
	sort.Strings(keys)
	for index, key := range keys {
		if index != 0 {
			sb.WriteRune('&')
		}
		sb.WriteString(key)
		sb.WriteRune('=')
		if value, found := parameters.Resolve(key); found {
			sb.WriteString(url.QueryEscape(value.String()))
		}

	}
	return sb.String()
}
