package main

import (
	"strings"

	elmo "github.com/okke/elmo/core"
)

// Module contains new elmo functions
//
var Module = elmo.NewModule("example", func(context elmo.RunContext) elmo.Value {
	return elmo.NewMappingForModule(context, []elmo.NamedValue{
		chipotle(),
		jalapeno()})
})

func chipotle() elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("chipotle", "", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		argLen, err := elmo.CheckArguments(arguments, 0, 1, "chipotle", "<string>")
		if err != nil {
			return err
		}

		if argLen == 0 {
			return elmo.NewStringLiteral("love them!")
		}

		value := elmo.EvalArgument(context, arguments[0])
		if value.Type() == elmo.TypeInteger {
			return elmo.NewStringLiteral(strings.Repeat("love them!", int(value.Internal().(int64))))
		}

		return elmo.NewErrorValue("please use nothing or an integer value as first argument")
	})
}

func jalapeno() elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("jalapeno", "", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		// Do it your self
		//
		return elmo.Nothing
	})
}
