package elmo

import "testing"

func dummyArgFromText(name string) Argument {
	return NewArgument(NewScriptMetaData("test", "test"), nil, NewIdentifier(name))
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

func TestLiteralStringsWithUnicodeCharacters(t *testing.T) {
	ParseTestAndRunBlock(t,
		`"⌘chipotle"`, ExpectValue(t, NewStringLiteral("⌘chipotle")))
}

func TestLiteralStringsWithBlock(t *testing.T) {
	ParseTestAndRunBlock(t,
		`"chipotle_\{}_jalapeno"`, ExpectValue(t, NewStringLiteral("chipotle__jalapeno")))
	ParseTestAndRunBlock(t,
		`"chipotle_\{true}_habanero"`, ExpectValue(t, NewStringLiteral("chipotle_true_habanero")))
	ParseTestAndRunBlock(t,
		`"chipotle_\{true}_\{false}_jalapeno"`, ExpectValue(t, NewStringLiteral("chipotle_true_false_jalapeno")))
	ParseTestAndRunBlock(t,
		`a:3; "chipotle_\{a}_jalapeno"`, ExpectValue(t, NewStringLiteral("chipotle_3_jalapeno")))
	ParseTestAndRunBlock(t,
		`a:33; b:45; "chipotle_\{a}_\{b}_jalapeno"`, ExpectValue(t, NewStringLiteral("chipotle_33_45_jalapeno")))
	ParseTestAndRunBlock(t,
		`pepper:"chipotle"; s:"\{pepper}"`, ExpectValue(t, NewStringLiteral("chipotle")))
}

func TestLongLiteralStringsWithBlock(t *testing.T) {
	ParseTestAndRunBlock(t,
		"`chipotle_`{}_jalapeno`", ExpectValue(t, NewStringLiteral("chipotle__jalapeno")))
	ParseTestAndRunBlock(t,
		"`chipotle_`{true}_jalapeno`", ExpectValue(t, NewStringLiteral("chipotle_true_jalapeno")))
	ParseTestAndRunBlock(t,
		"`chipotle_`{true}_`{false}_jalapeno`", ExpectValue(t, NewStringLiteral("chipotle_true_false_jalapeno")))
	ParseTestAndRunBlock(t,
		"``{true}`{false}`", ExpectValue(t, NewStringLiteral("truefalse")))
	ParseTestAndRunBlock(t,
		"a:3; `chipotle_`{a}_jalapeno`", ExpectValue(t, NewStringLiteral("chipotle_3_jalapeno")))
	ParseTestAndRunBlock(t,
		"a:33; b:45; `chipotle_`{a}_`{b}_jalapeno`", ExpectValue(t, NewStringLiteral("chipotle_33_45_jalapeno")))
	ParseTestAndRunBlock(t,
		"pepper:`chipotle`; s:``{pepper}`", ExpectValue(t, NewStringLiteral("chipotle")))
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

func TestStringAccessWithUTF8Strings(t *testing.T) {

	ParseTestAndRunBlock(t,
		`set s "😁😂😃"
		 s 0`, ExpectValue(t, NewStringLiteral("😁")))
	ParseTestAndRunBlock(t,
		`set s "😁😂😃"
		 s 1`, ExpectValue(t, NewStringLiteral("😂")))
	ParseTestAndRunBlock(t,
		`set s "😁😂😃"
		 s -1`, ExpectValue(t, NewStringLiteral("😃")))
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
		`c: {h: {i: {p: {o: {t: {l: {e: "sauce"}}}}}}}
		 c.h.i.p.o.t.l.e`, ExpectValue(t, NewStringLiteral("sauce")))

	ParseTestAndRunBlock(t,
		`c: {h: {i: {p: {o: {t: {l: {e: "sauce"}}}}}}}
 		 i: $c.h.i.p.o.t.l.e
		 i`, ExpectValue(t, NewStringLiteral("sauce")))

	ParseTestAndRunBlock(t,
		`c: {h: {i: {p: {o: {t: {l: {e: "sauce"}}}}}}}
  	 i: (c.h.i.p.o.t.l.e)
 		 i`, ExpectValue(t, NewStringLiteral("sauce")))
}

func TestPipeMultipleReturnValues(t *testing.T) {
	ParseTestAndRunBlock(t,
		`f: (func {return 3 4})
		 g: (func x y {return (multiply $x $y)})
		 f |g`, ExpectValue(t, NewIntegerLiteral(12)))
}

func TestPipeToList(t *testing.T) {
	ParseTestAndRunBlock(t,
		`f: (func {return 1})
		 l: [99 98 97]
		 f |l`, ExpectValue(t, NewIntegerLiteral(98)))
}

func TestPipeToString(t *testing.T) {
	ParseTestAndRunBlock(t,
		`f: (func {return 1})
		 s: "chipotle"
		 f |s`, ExpectValue(t, NewStringLiteral("h")))
}

func TestBlockWithoutCallsShouldReturnNothing(t *testing.T) {

	result := NewBlock(nil, nil, []Call{}).Run(NewRunContext(nil), []Argument{})

	if result != Nothing {
		t.Error("empty block should return nothing")
	}
}

func TestBlockWithOneCallsShouldReturnCallResult(t *testing.T) {

	context := NewRunContext(nil)

	context.Set("chipotle", NewStringLiteral("sauce"))

	result := NewBlock(nil, nil, []Call{NewCall(nil, nil, dummyArgFromText("chipotle"), []Argument{}, nil)}).Run(context, []Argument{})

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

	result := NewBlock(nil, nil, []Call{
		NewCall(nil, nil, dummyArgFromText("chipotle"), []Argument{}, nil),
		NewCall(nil, nil, dummyArgFromText("blackbeans"), []Argument{}, nil)}).Run(context, []Argument{})

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

	context.SetNamed(NewGoFunctionWithHelp("sauce", "", func(functionContext RunContext, arguments []Argument) Value {
		return NewStringLiteral("chipotle")
	}))

	result := NewBlock(nil, nil, []Call{NewCall(nil, nil, dummyArgFromText("sauce"), []Argument{}, nil)}).Run(context, []Argument{})

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

	context.SetNamed(NewGoFunctionWithHelp("echo", "", func(functionContext RunContext, arguments []Argument) Value {
		return arguments[0].Value()
	}))

	result := NewBlock(nil, nil, []Call{NewCall(nil, nil, dummyArgFromText("echo"), []Argument{NewArgument(nil, nil, NewStringLiteral("chipotle"))}, nil)}).Run(context, []Argument{})

	if result.String() != "chipotle" {
		t.Errorf("function should return (chipotle) instead of %s", result.String())
	}

}

func TestGoFunctionCanAlterContext(t *testing.T) {

	context := NewRunContext(nil)

	context.SetNamed(NewGoFunctionWithHelp("alter", "", func(functionContext RunContext, arguments []Argument) Value {
		context.Set(arguments[0].Value().String(), arguments[1].Value())
		return Nothing
	}))

	NewBlock(nil, nil, []Call{NewCall(nil, nil, dummyArgFromText("alter"), []Argument{
		NewArgument(nil, nil, NewStringLiteral("chipotle")),
		NewArgument(nil, nil, NewStringLiteral("sauce"))}, nil)}).Run(context, []Argument{})

	result, found := context.Get("chipotle")

	if !found {
		t.Error("function should manipulate context")
	} else {
		if result.String() != "sauce" {
			t.Errorf("function should return (sauce) instead of %s", result.String())
		}
	}

}

func TestIdentifierToBinary(t *testing.T) {

	if value := NewIdentifier("chipotle").(SerializableValue).ToBinary().ToRegular(); value.String() != "chipotle" {
		t.Errorf("expected chipotle, not \"%s\"", value.String())
	}

	if value := NewIdentifier("peppers.chipotle").(SerializableValue).ToBinary().ToRegular(); value.String() != "peppers.chipotle" {
		t.Errorf("expected chipotle, not \"%s\"", value.String())
	}

}

func TestStringToBinary(t *testing.T) {
	if value := NewStringLiteral("chipotle").(SerializableValue).ToBinary().ToRegular(); value.String() != "chipotle" {
		t.Errorf("expected chipotle, not \"%s\"", value.String())
	}
}

func TestIntToBinary(t *testing.T) {
	if value := NewIntegerLiteral(42).(SerializableValue).ToBinary().ToRegular(); value.String() != "42" {
		t.Errorf("expected 42, not \"%s\"", value.String())
	}
}

func TestFloatToBinary(t *testing.T) {
	if value := NewFloatLiteral(42.99).(SerializableValue).ToBinary().ToRegular(); value.String() != "42.99" {
		t.Errorf("expected 42.99, not \"%s\"", value.String())
	}
}

func TestBooleanToBinary(t *testing.T) {
	if value := True.(SerializableValue).ToBinary().ToRegular(); value.String() != "true" {
		t.Errorf("expected true, not \"%s\"", value.String())
	}
}
