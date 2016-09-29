package list

import (
	"math"

	"github.com/okke/elmo/core"
)

// ListModule contains functions that operate on lists
//
var Module = elmo.NewModule("list", initModule)

func initModule(context elmo.RunContext) elmo.Value {
	return elmo.NewMappingForModule(context, []elmo.NamedValue{
		_new(),
		_append(),
		prepend(),
		each(),
		_map(),
		filter()})
}

func _new() elmo.NamedValue {
	return elmo.NewGoFunction("new", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {
		return elmo.ListConstructor(context, arguments)
	})
}

func _append() elmo.NamedValue {
	return elmo.NewGoFunction("append", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		argLen, ok, err := elmo.CheckArguments(arguments, 2, math.MaxInt16, "append", "<list> <value>*")
		if !ok {
			return err
		}

		// first argument of a list function can be an identifier with the name of the list
		//
		list := elmo.EvalArgumentOrSolveIdentifier(context, arguments[0])

		if list.Type() != elmo.TypeList {
			return elmo.NewErrorValue("invalid call to append, expect at list as first argument: usage append <list> <value>*")
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

		argLen, ok, err := elmo.CheckArguments(arguments, 2, math.MaxInt16, "prepend", "<list> <value>*")
		if !ok {
			return err
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

func getValueIndexAndBlock(context elmo.RunContext, arguments []elmo.Argument) (elmo.Value, string, string, elmo.Argument, bool) {

	argLen := len(arguments)

	if argLen < 3 || argLen > 4 {
		return nil, "", "", nil, false
	}

	list := elmo.EvalArgumentOrSolveIdentifier(context, arguments[0])

	if list.Type() != elmo.TypeList {
		return nil, "", "", nil, false
	}

	valueName := elmo.EvalArgument2String(context, arguments[1])
	var indexName string
	if argLen == 4 {
		indexName = elmo.EvalArgument2String(context, arguments[2])
	}
	block := arguments[argLen-1]

	return list, valueName, indexName, block, true
}

func runInBlock(context elmo.RunContext, valueName string, value elmo.Value, indexName string, index int, block elmo.Argument) elmo.Value {
	context.Set(valueName, value)
	if indexName != "" {
		context.Set(indexName, elmo.NewIntegerLiteral(int64(index)))
	}
	return block.Value().(elmo.Block).Run(context, elmo.NoArguments)
}

func each() elmo.NamedValue {
	return elmo.NewGoFunction("each", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		list, valueName, indexName, block, valid := getValueIndexAndBlock(context, arguments)

		if !valid {
			return elmo.NewErrorValue("invalid call to each: usage each <list> <value identifier> <index identifier>? <block>")
		}

		var result elmo.Value

		subContext := context.CreateSubContext()
		for index, value := range list.Internal().([]elmo.Value) {
			result = runInBlock(subContext, valueName, value, indexName, index, block)
		}

		return result

	})
}

func _map() elmo.NamedValue {
	return elmo.NewGoFunction("map", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		list, valueName, indexName, block, valid := getValueIndexAndBlock(context, arguments)

		if !valid {
			return elmo.NewErrorValue("invalid call to map: usage map <list> <value identifier> <index identifier>? <block>")
		}

		oldValues := list.Internal().([]elmo.Value)
		l := len(oldValues)
		newValues := make([]elmo.Value, l, l)

		subContext := context.CreateSubContext()
		for index, value := range oldValues {
			newValues[index] = runInBlock(subContext, valueName, value, indexName, index, block)
		}

		return elmo.NewListValue(newValues)
	})
}

func filter() elmo.NamedValue {
	return elmo.NewGoFunction("filter", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		list, valueName, indexName, block, valid := getValueIndexAndBlock(context, arguments)

		if !valid {
			return elmo.NewErrorValue("invalid call to filter: usage filter <list> <value identifier> <index identifier>? <block>")
		}

		oldValues := list.Internal().([]elmo.Value)
		newValues := []elmo.Value{}

		subContext := context.CreateSubContext()
		for index, value := range oldValues {
			if runInBlock(subContext, valueName, value, indexName, index, block) == elmo.True {
				newValues = append(newValues, value)
			}
		}

		return elmo.NewListValue(newValues)
	})
}
