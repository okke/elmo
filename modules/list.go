package el

import "github.com/okke/elmo/core"

// ListModule contains functions that operate on lists
//
var Module = elmo.NewModule("el", initModule)

func initModule(context elmo.RunContext) elmo.Value {

	mapping := make(map[string]elmo.Value)

	all := []elmo.NamedValue{_append()}

	for _, v := range all {
		mapping[v.Name()] = v
	}

	return elmo.NewDictionaryValue(mapping)
}

func _append() elmo.NamedValue {
	return elmo.NewGoFunction("append", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		argLen := len(arguments)

		// first argument is a list, rest of the arguments are appended to the list
		if argLen < 2 {
			return elmo.NewErrorValue("invalid call to append, expect at least 2 parameters: usage append <list> <value> <value>?")
		}

		// first argument of a list function can be an identifier with the name of the list
		//
		list := elmo.EvalArgumentOrSolveIdentifier(context, arguments[0])

		if list.Type() != elmo.TypeList {
			return elmo.NewErrorValue("invalid call to append, expect at list as first argument: usage append <list> <value> <value>?")
		}

		internal := list.Internal().([]elmo.Value)

		for i := 1; i < argLen; i++ {
			internal = append(internal, elmo.EvalArgument(context, arguments[i]))
		}

		return elmo.NewListValue(internal)

	})
}
