package elmo

import "fmt"

func valueOrNil(value Value) Value {
	if value == nil {
		return Nothing
	}
	return value
}

// EvalArgument evaluates given argument
//
func EvalArgument(context RunContext, argument Argument) Value {

	if argument.Type() == TypeCall {
		return valueOrNil(argument.Value().(Runnable).Run(context, NoArguments))
	}

	if argument.Type() == TypeBlock {
		return valueOrNil(argument.Value().(Block).CopyWithinContext(context))
	}

	if argument.Type() == TypeString {
		return valueOrNil(argument.Value().(StringValue).ResolveBlocks(context))
	}

	return valueOrNil(argument.Value())

}

// EvalArgumentOrSolveIdentifier evaluates given argument
//
func EvalArgumentOrSolveIdentifier(context RunContext, argument Argument) Value {

	if argument.Type() == TypeIdentifier {
		value, found := context.Get(argument.String())
		if found {
			return value
		}
		return NewErrorValue(fmt.Sprintf("could not find %v", argument.String()))
	}

	return EvalArgument(context, argument)

}

// EvalArgumentWithBlock evaluates given argument and if argument is a block
// it will evaluate block content
//
func EvalArgumentWithBlock(context RunContext, argument Argument) Value {

	if argument.Type() == TypeCall || argument.Type() == TypeBlock {
		return argument.Value().(Runnable).Run(context, NoArguments)
	}

	return argument.Value()

}

// EvalArgument2String evaluates given argument and returns it String presentation
//
func EvalArgument2String(context RunContext, argument Argument) string {

	return EvalArgument(context, argument).String()

}
