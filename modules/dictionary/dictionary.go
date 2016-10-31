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
		keys(), knows(), dict(), new()})
}

func dict() elmo.NamedValue {
	return elmo.NewGoFunction("dict", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		argLen, ok, err := elmo.CheckArguments(arguments, 0, math.MaxInt16, "dict", "<list> | <block> | <value>*")
		if !ok {
			return err
		}

		mapping := make(map[string]elmo.Value)

		if argLen == 1 {
			evaluated := elmo.EvalArgument(context, arguments[0])
			if evaluated.Type() == elmo.TypeBlock {
				return elmo.NewDictionaryWithBlock(context, evaluated.(elmo.Block))
			}

			if evaluated.Type() != elmo.TypeList {
				return elmo.NewErrorValue(fmt.Sprintf("dict needs a list as argument. Can not create dictionary from %v", evaluated))
			}

			values := evaluated.Internal().([]elmo.Value)

			if (len(values) % 2) != 0 {
				return elmo.NewErrorValue("dict can not create a dictionary from an odd number of elements")
			}

			var key elmo.Value

			for i, val := range values {
				if i%2 == 0 {
					key = val
				} else {
					mapping[key.String()] = val
				}
			}

		} else {

			if (argLen % 2) != 0 {
				return elmo.NewErrorValue("dict can not create a dictionary from an odd number of elements")
			}

			var key elmo.Value

			for i, arg := range arguments {
				if i%2 == 0 {
					key = elmo.EvalArgument(context, arg)
				} else {
					mapping[key.String()] = elmo.EvalArgument(context, arg)
				}
			}
		}

		return elmo.NewDictionaryValue(nil, mapping)
	})
}

func new() elmo.NamedValue {
	return elmo.NewGoFunction("new", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		argLen, ok, err := elmo.CheckArguments(arguments, 1, 2, "new", "<value>*")
		if !ok {
			return err
		}

		parent := elmo.EvalArgument(context, arguments[0])

		if parent.Type() != elmo.TypeDictionary {
			return elmo.NewErrorValue(fmt.Sprintf("new expects a dictionary, not %s", parent.String()))
		}

		var mapping map[string]elmo.Value
		if argLen == 2 {

			dict := elmo.EvalArgument(context, arguments[1])

			if dict.Type() != elmo.TypeBlock {
				return elmo.NewErrorValue(fmt.Sprintf("new can not construct dictionary from %s", dict.String()))
			}

			dict = elmo.NewDictionaryWithBlock(context, dict.(elmo.Block))

			mapping = dict.Internal().(map[string]elmo.Value)
		} else {
			mapping = make(map[string]elmo.Value)
		}

		return elmo.NewDictionaryValue(parent, mapping)
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
