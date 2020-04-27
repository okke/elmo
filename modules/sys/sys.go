package sys

import elmo "github.com/okke/elmo/core"

// Module contains system functions
//
var Module = elmo.NewModule("sys", initModule)

func initModule(context elmo.RunContext) elmo.Value {
	return elmo.NewMappingForModule(context, []elmo.NamedValue{
		undefined(),
		_exec(),
		env(),
		setEnv(),
	})
}

// use elmo's undefined feature to convert all undefined function calls
// , within the the context of this module, to os executable commands
//
func undefined() elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("?", "", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {
		_, err := elmo.CheckArguments(arguments, 2, 2, "?", "<name> <args>")
		if err != nil {
			return err
		}

		name := elmo.EvalArgument2String(context, arguments[0])
		arglist := elmo.EvalArgument(context, arguments[1])
		if arglist.Type() != elmo.TypeList {
			return elmo.NewErrorValue("expected a list of arguments")
		}

		var pipeFrom Command
		var argValues = arglist.Internal().([]elmo.Value)

		if len(argValues) > 0 && argValues[0].IsType(typeInfoCommand) {
			// used in a pipe construction
			pipeFrom = argValues[0].Internal().(Command)
			argValues = argValues[1:]
		}
		strArgs := make([]string, len(argValues))
		for i, v := range argValues {
			strArgs[i] = v.String()
		}

		return NewCommandValue(pipeFrom, name, strArgs)
	})
}

func _exec() elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("exec", `Executes a command`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {
		_, err := elmo.CheckArguments(arguments, 1, 1, "exec", "<command>")
		if err != nil {
			return err
		}

		// first argument of an exec function can be an identifier with the name of the command
		//
		resolvedCommand := elmo.EvalArgumentOrSolveIdentifier(context, arguments[0])

		if !resolvedCommand.IsType(typeInfoCommand) {
			return elmo.NewErrorValue("invalid call to sys.exec, expected a command as first parameter. usage: exec command")
		}

		actualCommand := resolvedCommand.Internal().(Command)

		return actualCommand.Execute()
	})
}

func env() elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("env", `get all environment variables or one specified variable`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {
		argLen, err := elmo.CheckArguments(arguments, 0, 1, "env", "<value>?")
		if err != nil {
			return err
		}

		if argLen == 1 {
			value, found := allEnvironmentVariables.Resolve(elmo.EvalArgument2String(context, arguments[0]))
			if !found {
				return elmo.Nothing
			}
			return value
		}

		return allEnvironmentVariables
	})
}

func setEnv() elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("setEnv", `Sets an environment variable`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {
		_, err := elmo.CheckArguments(arguments, 2, 2, "setEnv", "<name> <value>")
		if err != nil {
			return err
		}

		name := elmo.EvalArgument2String(context, arguments[0])
		value := elmo.EvalArgument2String(context, arguments[1])

		setEnvVar(name, value)

		return elmo.NewStringLiteral(value)
	})
}
