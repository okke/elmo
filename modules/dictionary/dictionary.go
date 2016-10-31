package dict

import (
	"fmt"
	"sort"

	"github.com/okke/elmo/core"
)

// Module contains functions that operate on maps
//
var Module = elmo.NewModule("dict", initModule)

func initModule(context elmo.RunContext) elmo.Value {
	return elmo.NewMappingForModule(context, []elmo.NamedValue{
		keys(), knows(), new()})
}

func new() elmo.NamedValue {
	return elmo.NewGoFunction("new", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		argLen, ok, err := elmo.CheckArguments(arguments, 1, 2, "new", "<parent> <block|dictionary>")
		if !ok {
			return err
		}

		if argLen == 1 {

			// new without parent
			//
			dictValues := elmo.EvalArgument(context, arguments[0])

			if dictValues.Type() == elmo.TypeBlock {
				return elmo.NewDictionaryWithBlock(context, dictValues.(elmo.Block))
			}

			if dictValues.Type() == elmo.TypeDictionary {
				return elmo.NewDictionaryValue(nil, dictValues.Internal().(map[string]elmo.Value))
			}

			if dictValues.Type() == elmo.TypeList {
				return elmo.NewDictionaryFromList(nil, dictValues.Internal().([]elmo.Value))
			}

			return elmo.NewErrorValue(fmt.Sprintf("can not create dictionary from %v", dictValues))

		}

		parent := elmo.EvalArgumentOrSolveIdentifier(context, arguments[0])

		if parent.Type() != elmo.TypeDictionary {
			return elmo.NewErrorValue(fmt.Sprintf("invalid dictionary parent %v", parent))
		}

		dictValues := elmo.EvalArgument(context, arguments[1])

		if dictValues.Type() == elmo.TypeBlock {
			return elmo.NewDictionaryValue(parent, elmo.NewDictionaryWithBlock(context, dictValues.(elmo.Block)).Internal().(map[string]elmo.Value))
		}

		if dictValues.Type() == elmo.TypeDictionary {
			return elmo.NewDictionaryValue(parent, dictValues.Internal().(map[string]elmo.Value))
		}

		if dictValues.Type() == elmo.TypeList {
			return elmo.NewDictionaryFromList(parent, dictValues.Internal().([]elmo.Value))
		}

		return elmo.NewErrorValue(fmt.Sprintf("new can not construct dictionary from %s", dictValues.String()))

	})
}

func keys() elmo.NamedValue {
	return elmo.NewGoFunction("keys", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		_, ok, err := elmo.CheckArguments(arguments, 1, 1, "keys", "<dictionary>")
		if !ok {
			return err
		}

		// first argument of a dictionary function can be an identifier with the name of the dictionary
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

func knows() elmo.NamedValue {
	return elmo.NewGoFunction("knows", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		_, ok, err := elmo.CheckArguments(arguments, 2, 2, "knows", "<dictionary> <key>")
		if !ok {
			return err
		}

		// first argument of a dictionary function can be an identifier with the name of the dictionary
		//
		dict := elmo.EvalArgumentOrSolveIdentifier(context, arguments[0])

		if dict.Type() != elmo.TypeDictionary {
			return elmo.NewErrorValue("invalid call to keys, expect a dictionary as first argument: usage keys <dictionary>")
		}

		key := elmo.EvalArgument2String(context, arguments[1])

		mapping := dict.Internal().(map[string]elmo.Value)

		_, found := mapping[key]

		return elmo.NewBooleanLiteral(found)

	})
}
