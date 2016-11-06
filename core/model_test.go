package elmo

import "testing"

func dummyArgFromText(name string) Argument {
	return NewArgument(NewScriptMetaData("test", "test"), 0, 0, NewIdentifier(name))
}

func TestLiteralStrings(t *testing.T) {
	ParseTestAndRunBlock(t,
		`"chipotle\tjalapeno"`, ExpectValue(t, NewStringLiteral("chipotle\tjalapeno")))

	ParseTestAndRunBlock(t,
		`"chipotle\njalapeno"`, ExpectValue(t, NewStringLiteral("chipotle\njalapeno")))

	ParseTestAndRunBlock(t,
		"\"chipotle\\\\jalapeno\"", ExpectValue(t, NewStringLiteral("chipotle\\jalapeno")))

	ParseTestAndRunBlock(t,
		`"chipotle\-jalapeno"`, ExpectValue(t, NewStringLiteral("chipotle-jalapeno")))

	ParseTestAndRunBlock(t,
		"`chipotle\njalapeno`", ExpectValue(t, NewStringLiteral("chipotle\njalapeno")))

	ParseTestAndRunBlock(t,
		"`chipotle\\\\jalapeno`", ExpectValue(t, NewStringLiteral("chipotle\\\\jalapeno")))

	ParseTestAndRunBlock(t,
		"`chipotle'jalapeno`", ExpectValue(t, NewStringLiteral("chipotle'jalapeno")))

	ParseTestAndRunBlock(t,
		"`chipotle``jalapeno`", ExpectValue(t, NewStringLiteral("chipotle`jalapeno")))
}

func TestFunctionCallWithBlock(t *testing.T) {
	ParseTestAndRunBlock(t,
		`f: (func arg {return (type $arg)})
		 f {}`, ExpectValue(t, NewIdentifier("block")))
}

func TestMissingStatement(t *testing.T) {

	ParseTestAndRunBlock(t,
		`?: (func name args {
			return [$name $args]
		 })
		 chipotle uno dos tres`, ExpectValue(t, NewListValue([]Value{NewIdentifier("chipotle"),
			NewListValue([]Value{NewIdentifier("uno"), NewIdentifier("dos"), NewIdentifier("tres")})})))

	ParseTestAndRunBlock(t,
		`chipotle: {
			?: (func name args {
					return [$name $args]
				 })
		 }
		 chipotle.sauce uno dos tres`, ExpectValue(t, NewListValue([]Value{NewIdentifier("sauce"),
			NewListValue([]Value{NewIdentifier("uno"), NewIdentifier("dos"), NewIdentifier("tres")})})))

}

func TestStringAccess(t *testing.T) {

	ParseTestAndRunBlock(t,
		`set s "chipotle"
		 s 0`, ExpectValue(t, NewStringLiteral("c")))

	ParseTestAndRunBlock(t,
		`set s "chipotle"
			s 0 2`, ExpectValue(t, NewStringLiteral("chi")))

	ParseTestAndRunBlock(t,
		`set s "chipotle"
		 s 2 0`, ExpectValue(t, NewStringLiteral("ihc")))

	ParseTestAndRunBlock(t,
		`set s "chipotle"
		 s -1 0`, ExpectValue(t, NewStringLiteral("eltopihc")))

	ParseTestAndRunBlock(t,
		`set s "chipotle"
 		 s 1 2 3`, ExpectErrorValueAt(t, 2))

	ParseTestAndRunBlock(t,
		`set s "chipotle"
  	 s "sauce"`, ExpectErrorValueAt(t, 2))

	ParseTestAndRunBlock(t,
		`set s "chipotle"
	 	 s 1 "sauce"`, ExpectErrorValueAt(t, 2))

	ParseTestAndRunBlock(t,
		`set s "chipotle"
	 	 s "sauce" 2`, ExpectErrorValueAt(t, 2))

	ParseTestAndRunBlock(t,
		`set s "chipotle"
	 	 s 99`, ExpectErrorValueAt(t, 2))
}

func TestLiteralsAsCalls(t *testing.T) {

	ParseTestAndRunBlock(t,
		`[3]`, ExpectValue(t, NewListValue([]Value{NewIntegerLiteral(3)})))

	ParseTestAndRunBlock(t,
		`[3 4]`, ExpectValue(t, NewListValue([]Value{NewIntegerLiteral(3), NewIntegerLiteral(4)})))

	ParseTestAndRunBlock(t,
		`[3 "chipotle"]`, ExpectValue(t, NewListValue([]Value{NewIntegerLiteral(3), NewStringLiteral("chipotle")})))

	ParseTestAndRunBlock(t,
		`[[3 "chipotle"]]`, ExpectValue(t, NewListValue([]Value{NewListValue([]Value{NewIntegerLiteral(3), NewStringLiteral("chipotle")})})))

	ParseTestAndRunBlock(t,
		`3`, ExpectValue(t, NewIntegerLiteral(3)))

	ParseTestAndRunBlock(t,
		`3.0`, ExpectValue(t, NewFloatLiteral(3.0)))

	ParseTestAndRunBlock(t,
		`"chipotle"`, ExpectValue(t, NewStringLiteral("chipotle")))
}

func TestCallsAsCall(t *testing.T) {
	ParseTestAndRunBlock(t,
		`pepper: (func { return chipotle })
	   chipotle: (func x { return (plus $x 1) })
	   $pepper 2`, ExpectValue(t, NewIntegerLiteral(3)))

	ParseTestAndRunBlock(t,
		`pepper: (func { return jalapeno })
		 jalapeno: (func { return chipotle })
		 chipotle: (func x { return (plus $x 1) })
		 ($pepper) 2`, ExpectValue(t, NewIntegerLiteral(3)))
}

func TestDeepIdentifierLookup(t *testing.T) {
	ParseTestAndRunBlock(t,
		`m: {a: {b:3}}
		 m.a.b`, ExpectValue(t, NewIntegerLiteral(3)))

	ParseTestAndRunBlock(t,
		`m: {a: {b:3}}
 		 i: $m.a.b
		 i`, ExpectValue(t, NewIntegerLiteral(3)))

	ParseTestAndRunBlock(t,
		`m: {a: {b:3}}
  	 i: (m.a.b)
 		 i`, ExpectValue(t, NewIntegerLiteral(3)))
}

func TestBlockWithoutCallsShouldReturnNothing(t *testing.T) {

	result := NewBlock(nil, 0, 0, []Call{}).Run(NewRunContext(nil), []Argument{})

	if result != Nothing {
		t.Error("empty block should return nothing")
	}
}

func TestBlockWithOneCallsShouldReturnCallResult(t *testing.T) {

	context := NewRunContext(nil)

	context.Set("chipotle", NewStringLiteral("sauce"))

	result := NewBlock(nil, 0, 0, []Call{NewCall(nil, 0, 0, dummyArgFromText("chipotle"), []Argument{}, nil)}).Run(context, []Argument{})

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
		NewCall(nil, 0, 0, dummyArgFromText("chipotle"), []Argument{}, nil),
		NewCall(nil, 0, 0, dummyArgFromText("blackbeans"), []Argument{}, nil)}).Run(context, []Argument{})

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

	result := NewBlock(nil, 0, 0, []Call{NewCall(nil, 0, 0, dummyArgFromText("sauce"), []Argument{}, nil)}).Run(context, []Argument{})

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

	result := NewBlock(nil, 0, 0, []Call{NewCall(nil, 0, 0, dummyArgFromText("echo"), []Argument{NewArgument(nil, 0, 0, NewStringLiteral("chipotle"))}, nil)}).Run(context, []Argument{})

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

	NewBlock(nil, 0, 0, []Call{NewCall(nil, 0, 0, dummyArgFromText("alter"), []Argument{
		NewArgument(nil, 0, 0, NewStringLiteral("chipotle")),
		NewArgument(nil, 0, 0, NewStringLiteral("sauce"))}, nil)}).Run(context, []Argument{})

	result, found := context.Get("chipotle")

	if !found {
		t.Error("function should manipulate context")
	} else {
		if result.String() != "sauce" {
			t.Errorf("function should return (sauce) instead of %s", result.String())
		}
	}

}
