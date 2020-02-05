package elmo

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"plugin"
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

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (loader *loader) loadFromPlugin(folderName string, name string) Value {

	source := strings.Join([]string{folderName, "/", name, ".so"}, "")

	if !fileExists(source) {
		return nil
	}

	modulePlugin, err := plugin.Open(source)
	if err != nil {
		return NewErrorValue(err.Error())
	}

	moduleInitializer, err := modulePlugin.Lookup("ElmoPlugin")
	if err != nil {
		return NewErrorValue(err.Error())
	}

	moduleInitializerFunc, couldCast := moduleInitializer.(func(string) Module)
	if !couldCast {
		return NewErrorValue("found module initializer in shared library but it's not of type 'func(string) Module'")
	}

	module := moduleInitializerFunc(name)

	return module.Content(loader.context)
}

func (loader *loader) loadFromDir(folderName string, name string) Value {
	source := strings.Join([]string{folderName, "/", name, ".mo"}, "")

	b, err := ioutil.ReadFile(source)
	if err == nil {

		subContext := loader.context.CreateSubContext()

		result := ParseAndRunWithFile(subContext, string(b), source)

		if result.Type() == TypeError {
			result.(ErrorValue).Panic()
			return result
		}

		return NewDictionaryValue(nil, subContext.Mapping())
	}

	// when file exists but could not be read, return error
	//
	if fileExists(source) {
		return NewErrorValue(err.Error())
	}

	// could be a golang plugin
	//
	return loader.loadFromPlugin(folderName, name)
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

	err := NewErrorValue(fmt.Sprintf("could not find %s", name))
	return err.Panic()
}

// NewLoader constructs a new source code loader
//
func NewLoader(context RunContext, folders []string) Loader {
	return &loader{context: context, folders: folders}
}
