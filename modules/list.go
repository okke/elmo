package el

import "github.com/okke/elmo/core"

// ListModule contains functions that operate on lists
//
var Module = elmo.NewModule("el", initModule)

func initModule(context elmo.RunContext) elmo.Value {
	return elmo.NewMappingForModule(context, []elmo.NamedValue{
		_append(),
		prepend(),
		each()})
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

func prepend() elmo.NamedValue {
	return elmo.NewGoFunction("prepend", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		argLen := len(arguments)

		// first argument is a list, rest of the arguments are appended to the list
		if argLen < 2 {
			return elmo.NewErrorValue("invalid call to prepend, expect at least 2 parameters: usage prepend <list> <value> <value>?")
		}

		// first argument of a list function can be an identifier with the name of the list
		//
		list := elmo.EvalArgumentOrSolveIdentifier(context, arguments[0])

		if list.Type() != elmo.TypeList {
			return elmo.NewErrorValue("invalid call to prepend, expect at list as first argument: usage prepend <list> <value> <value>?")
		}

		internal := list.Internal().([]elmo.Value)

		for i := 1; i < argLen; i++ {
			internal = append([]elmo.Value{elmo.EvalArgument(context, arguments[i])}, internal...)
		}

		return elmo.NewListValue(internal)
	})
}

func each() elmo.NamedValue {
	return elmo.NewGoFunction("each", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		argLen := len(arguments)

		if argLen != 3 {
			return elmo.NewErrorValue("invalid call to each, expect at 3 parameters: usage each <list> <identifier> <block>")
		}

		// first argument of a list function can be an identifier with the name of the list
		//
		list := elmo.EvalArgumentOrSolveIdentifier(context, arguments[0])
		name := elmo.EvalArgument2String(context, arguments[1])
		block := arguments[argLen-1]

		if list.Type() != elmo.TypeList {
			return elmo.NewErrorValue("invalid call to each, expect at list as first argument: usage each <list> <identifier> <block>")
		}

		var result elmo.Value
		subContext := context.CreateSubContext()
		for _, v := range list.Internal().([]elmo.Value) {
			subContext.Set(name, v)
			result = block.Value().(elmo.Block).Run(subContext, elmo.NoArguments)
		}

		return result

	})
}
