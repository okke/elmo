package elmo

import "testing"

func TestBlockWithoutCallsShouldReturnNothing(t *testing.T) {

	result := NewBlock([]Call{}).Run(NewRunContext(nil))

	if result != Nothing {
		t.Error("empty block should return nothing")
	}
}

func TestBlockWithOneCallsShouldReturnCallResult(t *testing.T) {

	context := NewRunContext(nil)

	context.Set("chipotle", NewStringLiteral("sauce"))

	result := NewBlock([]Call{NewCall("chipotle", []Argument{})}).Run(context)

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

	result := NewBlock([]Call{
		NewCall("chipotle", []Argument{}),
		NewCall("blackbeans", []Argument{})}).Run(context)

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

	context.SetNamed(NewGoFunction("sauce", func(functionContext RunContext, values []Argument) Value {
		return NewStringLiteral("chipotle")
	}))

	result := NewBlock([]Call{NewCall("sauce", []Argument{})}).Run(context)

	if result == Nothing {
		t.Error("block with statement should return something")
	} else {
		if result.String() != "chipotle" {
			t.Errorf("block should return (chipotle) instead of %s", result.String())
		}
	}
}
