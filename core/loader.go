package elmo

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
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

func (loader *loader) buildFromGoCode(goModPath string) error {

	cmd := exec.Command("go", "build", "-buildmode=plugin")
	cmd.Dir = filepath.Dir(goModPath)
	_, err := cmd.Output()
	return errors.Wrap(err, fmt.Sprintf("could not exec: %s in %s", cmd.String(), cmd.Dir))

}

func (loader *loader) loadFromPlugin(folderName string, name string) Value {

	source := strings.Join([]string{folderName, "/", name, ".so"}, "")

	if !fileExists(source) {
		goModPath := strings.Join([]string{filepath.Dir(source), "/go.mod"}, "")
		if !fileExists(goModPath) {
			return nil
		}
		// found go source code, try to compile it
		//
		if err := loader.buildFromGoCode(goModPath); err != nil {
			return NewErrorValue(err.Error())
		}

		if !fileExists(source) {
			return NewErrorValue(fmt.Sprintf("found go code in %s but it did not compile to %s.so", folderName, name))
		}
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

func (loader *loader) addFileWatcher(source string, loaded DictionaryValue) {

	// TODO should we do this always or only when debugging?

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Op&fsnotify.Write == fsnotify.Write {

					reconstructed, _, _ := loader.constructDictionary(source)

					if reconstructed != nil {
						loaded.Replace(reconstructed)
					}
				}
			case _, ok := <-watcher.Errors:
				if !ok {
					return
				}
			}
		}
	}()

	err = watcher.Add(source)
	if err != nil {
		panic(err)
	}
}

func (loader *loader) constructDictionary(source string) (DictionaryValue, ErrorValue, error) {

	b, err := ioutil.ReadFile(source)
	if err != nil {
		return nil, nil, err
	}

	subContext := loader.context.CreateSubContext()

	result := ParseAndRunWithFile(subContext, string(b), source)

	if result.Type() == TypeError {
		result.(ErrorValue).Panic()
		return nil, result.(ErrorValue), nil
	}

	return NewDictionaryValue(nil, subContext.Mapping()), nil, nil

}

func (loader *loader) loadFromDir(folderName string, name string) Value {
	source := strings.Join([]string{folderName, "/", name, ".mo"}, "")

	constructed, loadError, _ := loader.constructDictionary(source)
	if constructed != nil {
		loader.addFileWatcher(source, constructed)
		return constructed
	}
	if loadError != nil {
		return loadError
	}

	// could be a golang plugin
	//
	return loader.loadFromPlugin(folderName, name)

}

func (loader *loader) Load(name string) Value {

	// try relative from current script
	//
	scriptName := loader.context.ScriptName()

	var result Value
	if scriptName != nil {
		result = loader.loadFromDir(path.Dir(scriptName.String()), name)
	} else {
		result = loader.loadFromDir(".", name)
	}
	if result != nil {
		return result
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
