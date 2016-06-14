package elmo

import "testing"

func TestSetAndGetValueFromContext(t *testing.T) {

	context := NewRunContext(nil)

	context.Set("sauce", NewStringLiteral("chipotle"))

	value, found := context.Get("sauce")

	if !found {
		t.Error("was expecting sauce to be in run context")
	} else {
		if value.String() != "chipotle" {
			t.Errorf("was expecting (chipotle) to be the sauce in run context, found (%s)", value.String())
		}
	}
}

func TestGetNonExistingValueFromContext(t *testing.T) {

	context := NewRunContext(nil)

	_, found := context.Get("sauce")

	if found {
		t.Error("was not expecting sauce to be in run context")
	}
}

func TestSetAndGetValueFromSubContext(t *testing.T) {

	context := NewRunContext(nil)

	context.Set("sauce", NewStringLiteral("chipotle"))

	sub := context.CreateSubContext()

	value, found := sub.Get("sauce")

	if !found {
		t.Error("was expecting sauce to be in run context")
	} else {
		if value.String() != "chipotle" {
			t.Errorf("was expecting (chipotle) to be the sauce in run context, found (%s)", value.String())
		}
	}
}

func TestGetNonExistingValueFromSubContext(t *testing.T) {

	context := NewRunContext(nil)

	sub := context.CreateSubContext()

	_, found := sub.Get("sauce")

	if found {
		t.Error("was not expecting sauce to be in run context")
	}
}

func TestSetAndGetNamedValue(t *testing.T) {

	context := NewRunContext(nil)

	context.SetNamed(NewGoFunction("sauce", func(values []Value) Value {
		return NewStringLiteral("chipotle")
	}))

	value, found := context.Get("sauce")

	if !found {
		t.Error("was expecting sauce to be in run context")
	} else {
		if value.Type() != TypeGoFunction {
			t.Errorf("was expecting sauce tp be a go function")
		}
	}
}
