package inspect

import elmo "github.com/okke/elmo/core"

// Module contains inspect functions
//
var Module = elmo.NewModule("inspect", initModule)

func initModule(context elmo.RunContext) elmo.Value {
	return elmo.NewMappingForModule(context, []elmo.NamedValue{
		meta()})
}

func meta() elmo.NamedValue {
	return elmo.NewGoFunction("meta", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		_, err := elmo.CheckArguments(arguments, 1, 1, "meta", "<inspectable>")
		if err != nil {
			return err
		}

		value := elmo.EvalArgument(context, arguments[0])
		inspectable, couldCast := value.(elmo.Inspectable)
		if !couldCast {
			return elmo.NewErrorValue("meta expects an inspectable value")
		}

		return elmo.NewDictionaryValue(nil, map[string]elmo.Value{
			"fileName": elmo.NewStringLiteral(inspectable.Meta().Name()),
			"beginsAt": elmo.NewIntegerLiteral(int64(inspectable.BeginsAt())),
			"length":   elmo.NewIntegerLiteral(int64(inspectable.EndsAt() - inspectable.BeginsAt())),
			"code": elmo.NewGoFunctionWithHelp("code", "get the actual elmo code", func(elmo.RunContext, []elmo.Argument) elmo.Value {
				content := inspectable.Meta().Content()
				return elmo.NewStringLiteral(string(content[int(inspectable.BeginsAt()):int(inspectable.EndsAt())]))
			})})

	})
}
