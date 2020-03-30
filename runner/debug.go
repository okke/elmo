package runner

import (
	"fmt"
	"strings"
	"time"

	elmo "github.com/okke/elmo/core"
)

func initDebugModule(runner Runner, debug bool) func(context elmo.RunContext) elmo.Value {
	return func(context elmo.RunContext) elmo.Value {
		return elmo.NewMappingForModule(context, []elmo.NamedValue{
			inDebug(debug),
			_log(debug),
			bp(runner, debug),
		})
	}
}

func doNothing(name string) elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp(name, "will do nothing",
		func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

			return elmo.Nothing
		})
}

func inDebug(debug bool) elmo.NamedValue {

	return elmo.NewGoFunctionWithHelp("inDebug", `check if elmo is running in debug mode`,
		func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {
			if debug {
				return elmo.True
			}
			return elmo.False
		})
}

func _log(debug bool) elmo.NamedValue {
	if !debug {
		return doNothing("log")
	}

	return elmo.NewGoFunctionWithHelp("log", `Log values to stdout
		Usage log <value>*`,
		func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {
			var builder strings.Builder

			builder.WriteString(time.Now().Format(time.RFC3339Nano))
			builder.WriteString(" ")
			for _, arg := range arguments {
				builder.WriteString(elmo.EvalArgument2String(context, arg))
			}

			fmt.Println(builder.String())

			return elmo.Nothing
		})
}

func bp(runner Runner, debug bool) elmo.NamedValue {

	if !debug {
		return doNothing("bp")
	}

	return elmo.NewGoFunctionWithHelp("bp", `Breakpoint into repl
		Usage bp <prompt prefix>?`,
		func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

			promptPrefix := "bp"

			if len(arguments) > 0 {
				promptPrefix = elmo.EvalArgument2String(context, arguments[0])
			}

			childRunner := runner.New(context.CreateSubContext(), promptPrefix)
			childRunner.Repl()

			return elmo.Nothing
		})
}
