package runner

import (
	"fmt"
	"strings"
	"time"

	elmo "github.com/okke/elmo/core"
)

func initDebugModule(debug bool) func(context elmo.RunContext) elmo.Value {
	return func(context elmo.RunContext) elmo.Value {
		return elmo.NewMappingForModule(context, []elmo.NamedValue{
			_log(debug),
			inDebug(debug),
		})
	}
}

func doNothing(name string) elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp(name, "will do nothing",
		func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

			return elmo.Nothing
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

func inDebug(debug bool) elmo.NamedValue {

	return elmo.NewGoFunctionWithHelp("inDebug", `check if elmo is running in debug mode`,
		func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {
			if debug {
				return elmo.True
			}
			return elmo.False
		})
}
