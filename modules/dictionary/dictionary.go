package dict

import (
	"fmt"
	"math"
	"sort"

	"github.com/okke/elmo/core"
)

// Module contains functions that operate on maps
//
var Module = elmo.NewModule("dict", initModule)

func initModule(context elmo.RunContext) elmo.Value {
	return elmo.NewMappingForModule(context, []elmo.NamedValue{
		new(), keys(), knows(), get(), merge(), set(), remove()})
}

func new() elmo.NamedValue {
	return elmo.NewGoFunction("new", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		argLen, err := elmo.CheckArguments(arguments, 1, 2, "new", "<parent> <block|dictionary>")
		if err != nil {
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

		_, err := elmo.CheckArguments(arguments, 1, 1, "keys", "<dictionary>")
		if err != nil {
			return err
		}

		// first argument of a dictionary function can be an identifier with the name of the dictionary
		//
		dict, ok := elmo.EvalArgumentOrSolveIdentifier(context, arguments[0]).(elmo.DictionaryValue)

		if !ok {
			return elmo.NewErrorValue("invalid call to keys, expect a dictionary as first argument: usage keys <dictionary>")
		}

		keyNames := dict.Keys()
		keys := make([]elmo.Value, len(keyNames))

		sort.Strings(keyNames)

		for i, k := range keyNames {
			keys[i] = elmo.NewIdentifier(k)
		}

		return elmo.NewListValue(keys)
	})
}

func knowsOrGet(name string, context elmo.RunContext, arguments []elmo.Argument) (elmo.Value, bool) {
	_, err := elmo.CheckArguments(arguments, 2, 2, name, "<dictionary> <key>")
	if err != nil {
		return err, false
	}

	// first argument of a dictionary function can be an identifier with the name of the dictionary
	//
	dict, ok := elmo.EvalArgumentOrSolveIdentifier(context, arguments[0]).(elmo.DictionaryValue)

	if !ok {
		return elmo.NewErrorValue(fmt.Sprintf("invalid call to %s, expect a dictionary as first argument", name)), false
	}

	key := elmo.EvalArgument(context, arguments[1])

	return dict.Resolve(key.String())

}

func knows() elmo.NamedValue {
	return elmo.NewGoFunction("knows", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		result, found := knowsOrGet("knows", context, arguments)
		if found {
			return elmo.True
		}

		// result can be an error
		if result.Type() == elmo.TypeError {
			return result
		}

		return elmo.False

	})
}

func get() elmo.NamedValue {
	return elmo.NewGoFunction("get", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		// get can act as knows by returning multiple return values
		//
		// set result found (d.dict $dict key)
		//

		result, found := knowsOrGet("get", context, arguments)
		if found {
			return elmo.NewReturnValue([]elmo.Value{result, elmo.True})
		}

		// result can be an error
		if result.Type() == elmo.TypeError {
			return result
		}

		return elmo.NewReturnValue([]elmo.Value{elmo.Nothing, elmo.False})

	})
}

func merge() elmo.NamedValue {
	return elmo.NewGoFunction("merge", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {
		argLen, err := elmo.CheckArguments(arguments, 2, math.MaxInt16, "merge", "<dictionary> <dictionary>+")
		if err != nil {
			return err
		}

		dictionaries := make([]elmo.DictionaryValue, argLen)
		for i, arg := range arguments {
			evaluated := elmo.EvalArgument(context, arg)

			dict, ok := evaluated.(elmo.DictionaryValue)
			if ok {
				dictionaries[i] = dict
			} else if arg.Type() == elmo.TypeBlock {
				dictionaries[i] = elmo.NewDictionaryWithBlock(context, evaluated.(elmo.Block)).(elmo.DictionaryValue)
			}
		}

		return dictionaries[0].Merge(dictionaries[1:])

	})
}

func set() elmo.NamedValue {
	return elmo.NewGoFunction("set!", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {
		_, err := elmo.CheckArguments(arguments, 3, 3, "set", "<dictionary> <symbol> <value>")
		if err != nil {
			return err
		}

		// first argument of a dictionary function can be an identifier with the name of the dictionary
		//
		dict, ok := elmo.EvalArgumentOrSolveIdentifier(context, arguments[0]).(elmo.DictionaryValue)

		if !ok {
			return elmo.NewErrorValue(fmt.Sprintf("invalid call to set!, expect a dictionary as first argument instead of %v", arguments[0]))
		}

		dict.Set(elmo.EvalArgument(context, arguments[1]), elmo.EvalArgument(context, arguments[2]))

		return dict.(elmo.Value)
	})
}

func remove() elmo.NamedValue {
	return elmo.NewGoFunction("remove!", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {
		_, err := elmo.CheckArguments(arguments, 2, 2, "remove", "<dictionary> <symbol>")
		if err != nil {
			return err
		}

		// first argument of a dictionary function can be an identifier with the name of the dictionary
		//
		dict, ok := elmo.EvalArgumentOrSolveIdentifier(context, arguments[0]).(elmo.DictionaryValue)

		if !ok {
			return elmo.NewErrorValue(fmt.Sprintf("invalid call to set!, expect a dictionary as first argument instead of %v", arguments[0]))
		}

		dict.Remove(elmo.EvalArgument(context, arguments[1]))

		return dict.(elmo.Value)
	})
}
