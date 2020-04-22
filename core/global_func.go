package elmo

import (
	"fmt"
	"math"
)

func createGoFunc(argNames []string, block Block) func(innerContext RunContext, innerArguments []Argument) Value {
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

		return block.Run(subContext, NoArguments)

	}
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

			argLen, err := CheckArguments(arguments, 1, math.MaxInt16, "func", "<identifier>* {...}")
			if err != nil {
				return err
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

			argNamesAsArgument := arguments[argStart : len(arguments)-1]
			block := EvalArgument(context, arguments[len(arguments)-1])
			argNames := make([]string, len(argNamesAsArgument))

			if argLen == 1 && block.Type() != TypeBlock {
				// block is not a block, maybe its an identifier that can be used
				// to lookup a function insted of creating on
				//
				result, found := context.Get(EvalArgument2String(context, arguments[len(arguments)-1]))
				if found && result.Type() == TypeGoFunction {
					return result
				}
			}

			for i, v := range argNamesAsArgument {
				argNames[i] = EvalArgument2String(context, v)
			}

			if block.Type() != TypeBlock {
				return NewErrorValue("invalid call to func, last parameter must be a block: usage func <identifier> <identifier>* {...}")
			}

			return NewGoFunctionWithBlock("anonymous", help, createGoFunc(argNames, block.(Block)), argNames, block.(Block))

		})
}
