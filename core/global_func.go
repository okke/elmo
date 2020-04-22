package elmo

import (
	"fmt"
)

func createGoFunc(argNames []string, evaluator func(evalContext RunContext) Value) func(innerContext RunContext, innerArguments []Argument) Value {
	return func(innerContext RunContext, innerArguments []Argument) Value {

		if len(argNames) != len(innerArguments) {
			return NewErrorValue(fmt.Sprintf("invalid call to user defined function: expect %d parameters instead of %d", len(argNames), len(innerArguments)))
		}

		cloneFrom := innerContext
		if cloneFrom.Parent() != nil {
			cloneFrom = cloneFrom.Parent()
		}
		subContext := cloneFrom.CreateSubContext()

		if innerContext.This() != nil {
			subContext.Set("this", innerContext.This())
		}

		for i, v := range innerArguments {
			subContext.Set(argNames[i], EvalArgument(innerContext, v))
		}

		return evaluator(subContext)

	}
}

func splitArgumentsForFunc(context RunContext, argStart int, arguments []Argument) ([]string, Value) {
	argNamesAsArgument := arguments[argStart : len(arguments)-1]
	code := EvalArgument(context, arguments[len(arguments)-1])
	argNames := make([]string, len(argNamesAsArgument))

	for i, v := range argNamesAsArgument {
		argNames[i] = EvalArgument2String(context, v)
	}
	return argNames, code
}

func _func() NamedValue {
	return NewGoFunctionWithHelp("func", `Create a new function
		Usage: func <help>? <symbol>* {...}
		Returns: a new function

		When first argument is a string, this value will be used as help text

		Given symbols denote function parameter names.

		Examples:

		> func a {}
		will create a function that accepts one parameter called 'a' (and does nothing)
		> func a { return $a }
		will create an echo function

		Note, function can be used once they are assigned to a variable

		> f: (func a { return $a })

		Example with help text:

		> f: (func a "we need more chipotles" {})
		> help f
		will result in "we need more chipotles"`,

		func(context RunContext, arguments []Argument) Value {

			argLen := len(arguments)
			if argLen == 0 {
				return NewErrorValue("func expects <help>? <identifier>* {...}")
			}

			argStart := 0

			help := ""
			if arguments[0].Type() == TypeString {
				if argLen == 1 {
					return NewErrorValue("func with help should at least have a body also")
				}
				help = EvalArgument2String(context, arguments[0])
				argStart = 1
			}

			argNames, code := splitArgumentsForFunc(context, argStart, arguments)

			if code.Type() != TypeBlock {
				return NewErrorValue("invalid call to func, last parameter must be a block")
			}

			block := code.(Block)

			evaluator := func(evalContext RunContext) Value {
				return block.Run(evalContext, NoArguments)
			}

			return NewGoFunctionWithBlock("anonymous", help, createGoFunc(argNames, evaluator), argNames, code.(Block))

		})
}

func template() NamedValue {
	return NewGoFunctionWithHelp("template", `Create a new template function
		Usage: template <help>? <symbol>* &"..."
		Returns: a new function

		When first argument is a string, this value will be used as help text

		Given symbols denote function parameter names.

		Examples:

		> echo: (template a &"\{$a}")
		will create a function that accepts one parameter actually returns it

		> goFileGen: (template packageName code &"package \{$packageName}\n\{$code}")
		will create a function which can generate a Go file
		`,

		func(context RunContext, arguments []Argument) Value {

			argLen := len(arguments)
			if argLen == 0 {
				return NewErrorValue("template expects <help>? <identifier>* &\"...\"")
			}

			argStart := 0

			help := ""
			if arguments[0].Type() == TypeString {
				if argLen > 1 {
					help = EvalArgument2String(context, arguments[0])
					argStart = 1
				}
			}

			argNames, code := splitArgumentsForFunc(context, argStart, arguments)

			if code.Type() != TypeString {
				return NewErrorValue("invalid call to template, last parameter must be a string")
			}

			str := code.(StringValue)

			evaluator := func(evalContext RunContext) Value {
				return str.ResolveBlocks(evalContext)
			}

			return NewGoFunctionWithHelp("template", help, createGoFunc(argNames, evaluator))

		})
}
