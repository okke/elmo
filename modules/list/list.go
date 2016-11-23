package list

import (
	"fmt"
	"math"

	"github.com/okke/elmo/core"
)

// Module contains functions that operate on lists
//
var Module = elmo.NewModule("list", initModule)

func initModule(context elmo.RunContext) elmo.Value {
	return elmo.NewMappingForModule(context, []elmo.NamedValue{
		_new(),
		_len(),
		at(),
		_append(),
		prepend(),
		each(),
		_map(),
		filter()})
}

func convertToList(value elmo.Value) ([]elmo.Value, elmo.ErrorValue) {

	if value.Type() == elmo.TypeList {
		return value.Internal().([]elmo.Value), nil
	}

	convertable, casted := value.Internal().(Listable)

	if !casted {
		return nil, elmo.NewErrorValue(fmt.Sprintf("can not convert %v to list", value))
	}

	return convertable.List(), nil

}

func _new() elmo.NamedValue {
	return elmo.NewGoFunction("new", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {
		return elmo.ListConstructor(context, arguments)
	})
}

func _len() elmo.NamedValue {
	return elmo.NewGoFunction("len", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		_, err := elmo.CheckArguments(arguments, 1, 1, "len", "<list>")
		if err != nil {
			return err
		}

		// first argument of a list function can be an identifier with the name of the list
		//
		list := elmo.EvalArgumentOrSolveIdentifier(context, arguments[0])

		internal, err := convertToList(list)
		if err != nil {
			return err
		}

		return elmo.NewIntegerLiteral(int64(len(internal)))
	})
}

func at() elmo.NamedValue {
	return elmo.NewGoFunction("at", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		_, err := elmo.CheckArguments(arguments, 2, 3, "at", "<list> <from> <to>?")
		if err != nil {
			return err
		}

		// first argument of a list function can be an identifier with the name of the list
		//
		list := elmo.EvalArgumentOrSolveIdentifier(context, arguments[0])

		internal, err := convertToList(list)
		if err != nil {
			return err
		}

		return elmo.NewListValue(internal).(elmo.Runnable).Run(context, arguments[1:])
	})
}

func _append() elmo.NamedValue {
	return elmo.NewGoFunction("append", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		argLen, err := elmo.CheckArguments(arguments, 2, math.MaxInt16, "append", "<list> <value>*")
		if err != nil {
			return err
		}

		// first argument of a list function can be an identifier with the name of the list
		//
		list := elmo.EvalArgumentOrSolveIdentifier(context, arguments[0])

		internal, err := convertToList(list)
		if err != nil {
			return err
		}

		for i := 1; i < argLen; i++ {
			internal = append(internal, elmo.EvalArgument(context, arguments[i]))
		}

		return elmo.NewListValue(internal)

	})
}

func prepend() elmo.NamedValue {
	return elmo.NewGoFunction("prepend", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		argLen, err := elmo.CheckArguments(arguments, 2, math.MaxInt16, "prepend", "<list> <value>*")
		if err != nil {
			return err
		}

		// first argument of a list function can be an identifier with the name of the list
		//
		list := elmo.EvalArgumentOrSolveIdentifier(context, arguments[0])

		internal, err := convertToList(list)
		if err != nil {
			return err
		}

		for i := 1; i < argLen; i++ {
			internal = append([]elmo.Value{elmo.EvalArgument(context, arguments[i])}, internal...)
		}

		return elmo.NewListValue(internal)
	})
}

func getValueIndexAndBlock(context elmo.RunContext, arguments []elmo.Argument) (elmo.Value, string, string, elmo.Value, bool) {

	argLen := len(arguments)

	if argLen < 2 || argLen > 4 {
		return nil, "", "", nil, false
	}

	list := elmo.EvalArgumentOrSolveIdentifier(context, arguments[0])

	if argLen == 2 {
		// eg. each list function_ref
		//
		return list, "", "", elmo.EvalArgument(context, arguments[argLen-1]), true
	}

	valueName := elmo.EvalArgument2String(context, arguments[1])
	var indexName string
	if argLen == 4 {
		indexName = elmo.EvalArgument2String(context, arguments[2])
	}
	block := elmo.EvalArgument(context, arguments[argLen-1])

	return list, valueName, indexName, block, true
}

func runInBlock(context elmo.RunContext, valueName string, value elmo.Value, indexName string, index int, block elmo.Value) elmo.Value {

	if block.Type() == elmo.TypeBlock {
		context.Set(valueName, value)
		if indexName != "" {
			context.Set(indexName, elmo.NewIntegerLiteral(int64(index)))
		}
		return block.(elmo.Block).Run(context, elmo.NoArguments)
	}

	if block.Type() == elmo.TypeGoFunction {
		return block.(elmo.Runnable).Run(context, []elmo.Argument{elmo.NewDynamicArgument(value)})
	}

	return elmo.NewErrorValue(fmt.Sprintf("invalid block %v", block))
}

func each() elmo.NamedValue {
	return elmo.NewGoFunction("each", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		list, valueName, indexName, block, valid := getValueIndexAndBlock(context, arguments)

		if !valid {
			return elmo.NewErrorValue("invalid call to each: usage each <list> <value identifier> <index identifier>? <block>")
		}

		var result elmo.Value

		subContext := context.CreateSubContext()

		internal, err := convertToList(list)
		if err != nil {
			return err
		}

		for index, value := range internal {
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

		oldValues, err := convertToList(list)
		if err != nil {
			return err
		}

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

		oldValues, err := convertToList(list)
		if err != nil {
			return err
		}

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
