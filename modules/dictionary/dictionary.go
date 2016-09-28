package dict

import (
	"sort"

	"github.com/okke/elmo/core"
)

// MapModule contains functions that operate on maps
//
var Module = elmo.NewModule("ed", initModule)

func initModule(context elmo.RunContext) elmo.Value {
	return elmo.NewMappingForModule(context, []elmo.NamedValue{
		keys()})
}

func keys() elmo.NamedValue {
	return elmo.NewGoFunction("keys", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {
		if len(arguments) != 1 {
			return elmo.NewErrorValue("invalid call to keys, expected exactly one parameter. usage: keys <dictionary>")
		}

		// first argument of a list function can be an identifier with the name of the list
		//
		dict := elmo.EvalArgumentOrSolveIdentifier(context, arguments[0])

		if dict.Type() != elmo.TypeDictionary {
			return elmo.NewErrorValue("invalid call to keys, expect a dictionary as first argument: usage keys <dictionary>")
		}

		mapping := dict.Internal().(map[string]elmo.Value)

		keyNames := make([]string, len(mapping))
		keys := make([]elmo.Value, len(mapping))

		i := 0
		for k := range mapping {
			keyNames[i] = k
			i++
		}

		sort.Strings(keyNames)

		for i, k := range keyNames {
			keys[i] = elmo.NewIdentifier(k)
		}

		return elmo.NewListValue(keys)
	})
}
