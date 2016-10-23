package str

import (
	"bytes"
	"math"

	elmo "github.com/okke/elmo/core"
)

// Module contains functions that operate on lists
//
var Module = elmo.NewModule("string", initModule)

func initModule(context elmo.RunContext) elmo.Value {
	return elmo.NewMappingForModule(context, []elmo.NamedValue{
		at(),
		_len(),
		join()})
}

func _len() elmo.NamedValue {
	return elmo.NewGoFunction("len", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		_, ok, err := elmo.CheckArguments(arguments, 1, 1, "len", "<string>")
		if !ok {
			return err
		}

		str := elmo.EvalArgument(context, arguments[0])
		if str.Type() == elmo.TypeError {
			return str
		}

		return elmo.NewIntegerLiteral(int64(len(str.String())))
	})
}

func at() elmo.NamedValue {
	return elmo.NewGoFunction("at", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		_, ok, err := elmo.CheckArguments(arguments, 2, 3, "at", "<string> <from> <to>?")
		if !ok {
			return err
		}

		str := elmo.EvalArgument(context, arguments[0])
		if str.Type() == elmo.TypeError {
			return str
		}

		return elmo.NewStringLiteral(str.String()).(elmo.Runnable).Run(context, arguments[1:])
	})
}

func join() elmo.NamedValue {
	return elmo.NewGoFunction("join", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		argLen, ok, err := elmo.CheckArguments(arguments, 2, math.MaxInt16, "join", "<string>*")
		if !ok {
			return err
		}

		var buffer bytes.Buffer

		for i := 0; i < argLen; i++ {
			value := elmo.EvalArgument(context, arguments[i])
			if value.Type() == elmo.TypeError {
				return value
			}
			buffer.WriteString(value.String())
		}

		return elmo.NewStringLiteral(buffer.String())

	})
}
