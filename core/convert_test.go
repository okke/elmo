package elmo

import "testing"

func TestConvertEmptyDictionary(t *testing.T) {

	goMap := ConvertDictionaryToMap(NewDictionaryValue(nil, map[string]Value{}))
	if len(goMap) != 0 {
		t.Errorf("expected an empty map")
	}
}

func TestConvertWithIgnoreType(t *testing.T) {

	if value := ConvertValueToInterface(
		NewGoFunction("f", func(RunContext, []Argument) Value { return nil }),
		TypeGoFunction); value != nil {

		t.Errorf("did not expect to convert a function, got %v", value)
	}
}

func TestConvertDictionaryWithValues(t *testing.T) {

	goMap := ConvertDictionaryToMap(NewDictionaryValue(nil, map[string]Value{
		"name": NewStringLiteral("chipotle"),
		"hotness": NewListValue([]Value{
			NewIntegerLiteral(30000),
			NewIntegerLiteral(50000),
		}),
		"likes": NewDictionaryValue(nil, map[string]Value{
			"me":  NewIntegerLiteral(5),
			"joe": NewIntegerLiteral(3),
		}),
	}))

	if len(goMap) != 3 {
		t.Errorf("expected a map with two keys")
	}

	if goMap["name"].(string) != "chipotle" {
		t.Errorf("expected a chipotle pepper")
	}

	if goMap["hotness"].([]interface{})[0].(int64) != 30000 {
		t.Errorf("expected a minimum sku of 30000")
	}

	if goMap["hotness"].([]interface{})[1].(int64) != 50000 {
		t.Errorf("expected a maximum sku of 50000")
	}

	if goMap["likes"].(map[string]interface{})["me"].(int64) != 5 {
		t.Errorf("well actualy I love chipotles so I expected 5 stars")
	}

	if goMap["likes"].(map[string]interface{})["joe"].(int64) != 3 {
		t.Errorf("joe only gave 3 stars")
	}
}

func TestConvertDictionaryWithFunction(t *testing.T) {

	goMap := ConvertDictionaryToMap(NewDictionaryValue(nil, map[string]Value{
		"f": NewGoFunction("f", func(RunContext, []Argument) Value { return nil }),
	}), TypeGoFunction)
	if len(goMap) != 0 {
		t.Errorf("expected an empty map, got %v", goMap["f"])
	}
}
