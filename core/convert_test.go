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
		NewGoFunctionWithHelp("f", "", func(RunContext, []Argument) Value { return nil }),
		TypeGoFunction); value != nil {

		t.Errorf("did not expect to convert a function, got %v", value)
	}
}

func TestConvertDictionaryWithValues(t *testing.T) {

	goMap := ConvertDictionaryToMap(NewDictionaryValue(nil, map[string]Value{
		"name": NewStringLiteral("chipotle"),
		"hotness": NewListValue([]Value{
			NewIntegerLiteral(3000),
			NewIntegerLiteral(5000),
		}),
		"likes": NewDictionaryValue(nil, map[string]Value{
			"me":  NewIntegerLiteral(5),
			"joe": NewIntegerLiteral(3),
		}),
	}))

	if len(goMap) != 3 {
		t.Errorf("expected a map with 3 keys")
	}

	if goMap["name"].(string) != "chipotle" {
		t.Errorf("expected a chipotle pepper")
	}

	if goMap["hotness"].([]interface{})[0].(int64) != 3000 {
		t.Errorf("expected a minimum sku of 30000")
	}

	if goMap["hotness"].([]interface{})[1].(int64) != 5000 {
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
		"f": NewGoFunctionWithHelp("f", "", func(RunContext, []Argument) Value { return nil }),
	}), TypeGoFunction)
	if len(goMap) != 0 {
		t.Errorf("expected an empty map, got %v", goMap["f"])
	}
}

func TestConvertDictionaryToFlatDictionary(t *testing.T) {

	goMap := ConvertDictionaryToMap(ConvertDictionaryToFlatDictionary(NewDictionaryValue(nil, map[string]Value{
		"likes": NewDictionaryValue(nil, map[string]Value{
			"me": NewIntegerLiteral(5),
			"joe": NewDictionaryValue(nil, map[string]Value{
				"badMood":  NewStringLiteral("arghh"),
				"goodMood": NewIntegerLiteral(3),
				"mixedMood": NewListValue([]Value{
					NewIntegerLiteral(3), NewIntegerLiteral(4), NewIntegerLiteral(5),
				}),
			}),
		}),
	})))

	if len(goMap) != 6 {
		t.Errorf("expected a map with 3 keys, not %v", goMap)
	}

	if i, found := goMap["likes_me"]; !found || i.(int64) != 5 {
		t.Errorf("well told you before, chipotles are 5 star peppers, check %v", goMap)
	}

	if i, found := goMap["likes_joe_badMood"]; !found || i.(string) != "arghh" {
		t.Errorf("when joe is in a bad mood, he hates chipotles, check %v", goMap)
	}

	if i, found := goMap["likes_joe_goodMood"]; !found || i.(int64) != 3 {
		t.Errorf("well told you before, joe love them not that much but still okay, check %v", goMap)
	}

	if i, found := goMap["likes_joe_mixedMood_0"]; !found || i.(int64) != 3 {
		t.Errorf("well told you before, joe love them not that much but still okay when he's got mixed feelings, check %v", goMap)
	}

	if i, found := goMap["likes_joe_mixedMood_1"]; !found || i.(int64) != 4 {
		t.Errorf("joe love them not when he's got mixed feelings, check %v", goMap)
	}

	if i, found := goMap["likes_joe_mixedMood_2"]; !found || i.(int64) != 5 {
		t.Errorf("joe really goes mad on chipotles when he's got mixed feelings, check %v", goMap)
	}
}

func TestConvertDictionaryToFlatDictionaryWithDictionariesInLists(t *testing.T) {

	goMap := ConvertDictionaryToMap(ConvertDictionaryToFlatDictionary(NewDictionaryValue(nil, map[string]Value{
		"likes": NewDictionaryValue(nil, map[string]Value{
			"me": NewListValue([]Value{
				NewDictionaryValue(nil, map[string]Value{
					"chipotle": NewIntegerLiteral(1),
					"jalapeno": NewIntegerLiteral(2),
				}),
				NewDictionaryValue(nil, map[string]Value{
					"chipotle": NewIntegerLiteral(3),
					"jalapeno": NewIntegerLiteral(4),
				}),
			}),
		}),
	})))

	if len(goMap) != 4 {
		t.Errorf("expected a map with 3 keys, not %v", goMap)
	}

	if i, found := goMap["likes_me_0_chipotle"]; !found || i.(int64) != 1 {
		t.Errorf("chipotle low score should be 1, check %v", goMap)
	}
	if i, found := goMap["likes_me_1_chipotle"]; !found || i.(int64) != 3 {
		t.Errorf("chipotle high score should be 3, check %v", goMap)
	}
	if i, found := goMap["likes_me_0_jalapeno"]; !found || i.(int64) != 2 {
		t.Errorf("jalapeno low score should be 2, check %v", goMap)
	}
	if i, found := goMap["likes_me_1_jalapeno"]; !found || i.(int64) != 4 {
		t.Errorf("jalapeno high score should be 4, check %v", goMap)
	}
}
