package inspect

import (
	"fmt"

	elmo "github.com/okke/elmo/core"
)

// Module contains inspect functions
//
var Module = elmo.NewModule("inspect", initModule)

func initModule(context elmo.RunContext) elmo.Value {
	return elmo.NewMappingForModule(context, []elmo.NamedValue{
		meta(), calls(), block()})
}

func meta() elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("meta", `
	usage: meta <block>
	Returns a dictionary with meta data of an inspectable element (block, call or argument)
	The returned dictionary looks like:
	{
		fileName: name of the elmo file/source this element declared
		beginsAt: absolute charactyer position in the file/source
		length: number of characters
		code: parsed code
	}

	a simple example is just by passing an empty block of code to it:

	> inspect: (load inspect)
	> meta: (inspect.meta {})

	`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		_, err := elmo.CheckArguments(arguments, 1, 1, "meta", "<inspectable>")
		if err != nil {
			return err
		}

		value := elmo.EvalArgument(context, arguments[0])
		inspectable, couldCast := value.(elmo.Inspectable)
		if !couldCast {
			return elmo.NewErrorValue(fmt.Sprintf("meta expects an inspectable value, not a value of type %v", value.Info().Name()))
		}

		dict := elmo.NewDictionaryValue(nil, map[string]elmo.Value{
			"fileName": elmo.NewStringLiteral(inspectable.Meta().Name()),
			"beginsAt": elmo.NewIntegerLiteral(int64(inspectable.BeginsAt())),
			"length":   elmo.NewIntegerLiteral(int64(inspectable.EndsAt() - inspectable.BeginsAt())),
			"code": elmo.NewGoFunctionWithHelp("code", "get the actual elmo code", func(elmo.RunContext, []elmo.Argument) elmo.Value {
				content := inspectable.Meta().Content()
				return elmo.NewStringLiteral(string(content[int(inspectable.BeginsAt()):int(inspectable.EndsAt())]))
			})})

		inspectable.Enrich(dict.(elmo.DictionaryValue))

		return dict

	})
}

func calls() elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("calls", `
	usage: calls <block>
	return a list of calls that are declared in the given block of code
	
	`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {
		_, err := elmo.CheckArguments(arguments, 1, 1, "calls", "<block>")
		if err != nil {
			return err
		}

		value := elmo.EvalArgument(context, arguments[0])
		if value.Type() != elmo.TypeBlock {
			return elmo.NewErrorValue("calls expects a block")
		}

		block := value.(elmo.Block)
		calls := block.Calls()

		values := make([]elmo.Value, len(calls), len(calls))

		for i, call := range calls {
			values[i] = call
		}

		return elmo.NewListValue(values)
	})
}

func block() elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("block", `
	usage: block <function>
	return the code block of a function
	
	`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {
		_, err := elmo.CheckArguments(arguments, 1, 1, "block", "<function>")
		if err != nil {
			return err
		}

		value := elmo.EvalArgument(context, arguments[0])

		userFunction, couldCast := value.(elmo.UserDefinedFunction)
		if !couldCast {
			return elmo.NewErrorValue("block expects a function")
		}

		block := userFunction.Block()
		if block == nil {
			return elmo.NewErrorValue("block expects a function written in elmo")
		}

		if block.Type() != elmo.TypeBlock {
			return elmo.NewErrorValue(fmt.Sprintf("block is not of type block, found %v", block))
		}

		return block
	})
}
