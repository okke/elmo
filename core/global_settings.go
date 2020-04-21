package elmo

import "sync"

// GlobalSettingData holds all global settings
//
type GlobalSettingData struct {
	Debug     bool
	HotReload bool
}

var createGlobalSettingsOnce sync.Once
var globalSettingSingleton *GlobalSettingData
var globalSettingsSingletonDictionary DictionaryValue

func createGlobalSettingSingletons() {
	createGlobalSettingsOnce.Do(func() {

		globalSettingSingleton = &GlobalSettingData{}
		globalSettingsSingletonDictionary = NewDictionaryFromStruct(nil, globalSettingSingleton)

	})
}

// GlobalSettings returns a singleton settings structure
//
func GlobalSettings() *GlobalSettingData {

	createGlobalSettingSingletons()

	return globalSettingSingleton
}

// GlobalSettingsDictionary returns a singleton dictionary with global settings
//
func GlobalSettingsDictionary() DictionaryValue {

	createGlobalSettingSingletons()

	return globalSettingsSingletonDictionary
}

func globalSettings() NamedValue {

	return NewGoFunctionWithHelp("globalSettings", `returns a dictionary with global elmo settings`,
		func(context RunContext, arguments []Argument) Value {
			return GlobalSettingsDictionary()
		})
}
