package elmo

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
)

type loader struct {
	context RunContext
	folders []string
}

// Loader is responsible for loading external elmo sources
//
type Loader interface {
	Load(name string) Value
}

func (loader *loader) loadFromDir(folderName string, name string) Value {
	source := strings.Join([]string{folderName, "/", name, ".mo"}, "")

	b, err := ioutil.ReadFile(source)
	if err == nil {

		subContext := loader.context.CreateSubContext()

		result := ParseAndRunWithFile(subContext, string(b), source)

		if result.Type() == TypeError {
			return result
		}

		return NewDictionaryValue(nil, subContext.Mapping())
	}
	return nil
}

func (loader *loader) Load(name string) Value {

	// try relative from current script
	//
	scriptName := loader.context.ScriptName()
	if scriptName != nil {
		result := loader.loadFromDir(path.Dir(scriptName.String()), name)
		if result != nil {
			return result
		}
	}

	// try known folders
	//
	for _, folderName := range loader.folders {
		result := loader.loadFromDir(folderName, name)
		if result != nil {
			return result
		}
	}

	return NewErrorValue(fmt.Sprintf("could not find %s", name))

}

// NewLoader constructs a new source code loader
//
func NewLoader(context RunContext, folders []string) Loader {
	return &loader{context: context, folders: folders}
}
