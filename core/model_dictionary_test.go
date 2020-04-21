package elmo

import (
	"testing"
)

type testStruct2DictionaryData struct {
	BoolField   bool
	StringField string
	IntField    int64
	FloatField  float64
}

func testField(t *testing.T, context RunContext, dict DictionaryValue, fieldName string, fieldType Type, getValue Value, setValue Value) {
	f, found := dict.Resolve(fieldName)
	if !found {
		t.Errorf("can not find field accessor %s", fieldName)
	}

	value := f.(Runnable).Run(context, []Argument{})
	if value.Type() != fieldType {
		t.Errorf("expect field accessor get to return %v", fieldType)
	}

	if fieldType == TypeString {
		if value.String() != getValue.String() {
			t.Errorf("expect field to be %v not %v", getValue, value)
		}
	} else {
		if value.Internal() != getValue.Internal() {
			t.Errorf("expect field to be %v not %v", getValue, value.Internal())
		}
	}

	value = f.(Runnable).Run(context, []Argument{NewDynamicArgument(setValue)})
	if value.Type() != fieldType {
		t.Errorf("expect field accessor set to return %v", fieldType)
	}

	if fieldType == TypeString {
		if value.String() != setValue.String() {
			t.Errorf("expect field after set to be %v not %v", setValue, value.Internal())
		}
	} else {
		if value.Internal() != setValue.Internal() {
			t.Errorf("expect field after set to be %v not %v", setValue, value.Internal())
		}
	}
}

func TestStruct2Dictionary(t *testing.T) {

	context := NewRunContext(nil)
	dict := NewDictionaryFromStruct(nil, &testStruct2DictionaryData{
		BoolField:   true,
		StringField: "chipotle",
		IntField:    42,
		FloatField:  42.24,
	})

	testField(t, context, dict, "BoolField", TypeBoolean, True, False)
	testField(t, context, dict, "StringField", TypeString, NewStringLiteral("chipotle"), NewStringLiteral("jalapeno"))
	testField(t, context, dict, "IntField", TypeInteger, NewIntegerLiteral(42), NewIntegerLiteral(24))
	testField(t, context, dict, "FloatField", TypeFloat, NewFloatLiteral(42.24), NewFloatLiteral(24.42))
}
