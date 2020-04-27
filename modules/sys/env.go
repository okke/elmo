package sys

import (
	"os"
	"strings"

	elmo "github.com/okke/elmo/core"
)

var allEnvironmentVariables = constructEnvDictionary()

func constructEnvDictionary() elmo.DictionaryValue {

	mapping := make(map[string]elmo.Value)

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		name := pair[0]

		mapping[name] = elmo.NewStringLiteral(os.Getenv(name))
	}

	return elmo.NewDictionaryValue(nil, mapping)
}

func setEnvVar(name, value string) {
	os.Setenv(name, value)
	allEnvironmentVariables.Set(elmo.NewStringLiteral(name), elmo.NewStringLiteral(os.Getenv(name)))
}
