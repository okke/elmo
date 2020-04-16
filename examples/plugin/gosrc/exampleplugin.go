package main

import (
	elmo "github.com/okke/elmo/core"
)

// ElmoPlugin is the entyry point for every elmo plugin that gets loaded
// through elmo's buildin load function
//
func ElmoPlugin(name string) elmo.Module {
	return elmo.NewModule(name, func(context elmo.RunContext) elmo.Value {
		return elmo.NewMappingForModule(context, []elmo.NamedValue{
			dosomething()})
	})
}

func dosomething() elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("forthytwo", `the answer to everything
		Usage: forthytwo
		Returns: 42`,

		func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

			_, err := elmo.CheckArguments(arguments, 0, 0, "dosomething", "")
			if err != nil {
				return err
			}

			return elmo.NewIntegerLiteral(42)
		})
}
