package elmo

import "testing"

func TestBlockWithoutCallsShouldReturnNothing(t *testing.T) {

	result := NewBlock(nil, 0, 0, []Call{}).Run(NewRunContext(nil), []Argument{})

	if result != Nothing {
		t.Error("empty block should return nothing")
	}
}

func TestBlockWithOneCallsShouldReturnCallResult(t *testing.T) {

	context := NewRunContext(nil)

	context.Set("chipotle", NewStringLiteral("sauce"))

	result := NewBlock(nil, 0, 0, []Call{NewCall(nil, 0, 0, []string{"chipotle"}, []Argument{})}).Run(context, []Argument{})

	if result == Nothing {
		t.Error("block with statement should return something")
	} else {
		if result.String() != "sauce" {
			t.Errorf("block should return (sauce) instead of %s", result.String())
		}
	}
}

func TestBlockWithTwoCallsShouldReturnLastCallResult(t *testing.T) {

	context := NewRunContext(nil)

	context.Set("chipotle", NewStringLiteral("sauce"))
	context.Set("blackbeans", NewStringLiteral("soup"))

	result := NewBlock(nil, 0, 0, []Call{
		NewCall(nil, 0, 0, []string{"chipotle"}, []Argument{}),
		NewCall(nil, 0, 0, []string{"blackbeans"}, []Argument{})}).Run(context, []Argument{})

	if result == Nothing {
		t.Error("block with statement should return something")
	} else {
		if result.String() != "soup" {
			t.Errorf("block should return (soup) instead of %s", result.String())
		}
	}
}

func TestBlockCallToNativeFunctionShouldExecuteFunction(t *testing.T) {

	context := NewRunContext(nil)

	context.SetNamed(NewGoFunction("sauce", func(functionContext RunContext, arguments []Argument) Value {
		return NewStringLiteral("chipotle")
	}))

	result := NewBlock(nil, 0, 0, []Call{NewCall(nil, 0, 0, []string{"sauce"}, []Argument{})}).Run(context, []Argument{})

	if result == Nothing {
		t.Error("block with statement should return something")
	} else {
		if result.String() != "chipotle" {
			t.Errorf("block should return (chipotle) instead of %s", result.String())
		}
	}
}

func TestGoFunctionWithOneArgumentCanReturnArgumentValue(t *testing.T) {

	context := NewRunContext(nil)

	context.SetNamed(NewGoFunction("echo", func(functionContext RunContext, arguments []Argument) Value {
		return arguments[0].Value()
	}))

	result := NewBlock(nil, 0, 0, []Call{NewCall(nil, 0, 0, []string{"echo"}, []Argument{NewArgument(nil, 0, 0, NewStringLiteral("chipotle"))})}).Run(context, []Argument{})

	if result.String() != "chipotle" {
		t.Errorf("function should return (chipotle) instead of %s", result.String())
	}

}

func TestGoFunctionCanAlterContext(t *testing.T) {

	context := NewRunContext(nil)

	context.SetNamed(NewGoFunction("alter", func(functionContext RunContext, arguments []Argument) Value {
		context.Set(arguments[0].Value().String(), arguments[1].Value())
		return Nothing
	}))

	NewBlock(nil, 0, 0, []Call{NewCall(nil, 0, 0, []string{"alter"}, []Argument{
		NewArgument(nil, 0, 0, NewStringLiteral("chipotle")),
		NewArgument(nil, 0, 0, NewStringLiteral("sauce"))})}).Run(context, []Argument{})

	result, found := context.Get("chipotle")

	if !found {
		t.Error("function should manipulate context")
	} else {
		if result.String() != "sauce" {
			t.Errorf("function should return (sauce) instead of %s", result.String())
		}
	}

}
