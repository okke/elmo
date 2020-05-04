package list

import (
	"fmt"
	"math"
	"sort"

	elmo "github.com/okke/elmo/core"
)

// Module contains functions that operate on lists
//
var Module = elmo.NewModule("list", initModule)

func initModule(context elmo.RunContext) elmo.Value {
	return elmo.NewMappingForModule(context, []elmo.NamedValue{
		_new(),
		tuple(),
		at(),
		_append(),
		mutableAppend(),
		prepend(),
		mutablePrepend(),
		each(),
		_map(),
		where(),
		_sort(),
		mutableSort()})
}

func convertToList(value elmo.Value) ([]elmo.Value, elmo.ErrorValue) {

	if value.Type() == elmo.TypeList {
		return value.Internal().([]elmo.Value), nil
	}

	convertable, casted := value.Internal().(elmo.Listable)

	if !casted {
		return nil, elmo.NewErrorValue(fmt.Sprintf("can not convert %v to list", value))
	}

	return convertable.List(), nil

}

func _new() elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("new", `create a new list`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {
		return elmo.ListConstructor(context, arguments)
	})
}

func tuple() elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("tuple", `convert a list to a tuple which can be piped into a function`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		_, err := elmo.CheckArguments(arguments, 1, 1, "tuple", "<list>")
		if err != nil {
			return err
		}

		list := elmo.EvalArgumentOrSolveIdentifier(context, arguments[0])

		internal, err := convertToList(list)
		if err != nil {
			return err
		}

		return elmo.NewReturnValue(internal)
	})
}

func at() elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("at", "acces list content", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

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

func appendAndOptionallyChange(name string, change bool) elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp(name, "append items to a list", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

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

		if change {
			mutable, ok := list.(elmo.MutableValue)
			if !ok {
				return elmo.NewErrorValue(fmt.Sprintf("can not mutate %v", list))
			}
			if _, err := mutable.Mutate(internal); err != nil {
				return err
			}
		}

		return elmo.NewListValue(internal)

	})
}

func _append() elmo.NamedValue {
	return appendAndOptionallyChange("append", false)
}

func mutableAppend() elmo.NamedValue {
	return appendAndOptionallyChange("append!", true)
}

func prependAndOptionallyChange(name string, change bool) elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp(name, "prepend items to a list", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {
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

		if change {
			mutable, ok := list.(elmo.MutableValue)
			if !ok {
				return elmo.NewErrorValue(fmt.Sprintf("can not mutate %v", list))
			}
			if _, err := mutable.Mutate(internal); err != nil {
				return err
			}
		}

		return elmo.NewListValue(internal)
	})
}

func prepend() elmo.NamedValue {
	return prependAndOptionallyChange("prepend", false)
}

func mutablePrepend() elmo.NamedValue {
	return prependAndOptionallyChange("prepend!", true)
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
	return elmo.NewGoFunctionWithHelp("each", `iterate over items in list and executes block of code`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

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
	return elmo.NewGoFunctionWithHelp("map", `map all items to values returned by block of code`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

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

func where() elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("where", `filter items`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		list, valueName, indexName, block, valid := getValueIndexAndBlock(context, arguments)

		if !valid {
			return elmo.NewErrorValue("invalid call to where: usage where <list> <value identifier> <index identifier>? <block>")
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

func mutableSort() elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("sort!", `sort items in list`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		_, err := elmo.CheckArguments(arguments, 1, 1, "sort!", "<list>")
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

		sort.Slice(internal, func(i, j int) bool {
			c, _ := internal[i].(elmo.ComparableValue).Compare(context, internal[j])
			return c < 0
		})

		return elmo.NewListValue(internal)
	})
}

func _sort() elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("sort", `created a sorted list`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		_, err := elmo.CheckArguments(arguments, 1, 1, "sort", "<list>")
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

		copyOfInternal := make([]elmo.Value, len(internal), len(internal))
		copy(copyOfInternal, internal)

		sort.Slice(copyOfInternal, func(i, j int) bool {
			c, _ := copyOfInternal[i].(elmo.ComparableValue).Compare(context, copyOfInternal[j])
			return c < 0
		})

		return elmo.NewListValue(copyOfInternal)
	})
}
