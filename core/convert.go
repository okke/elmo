package elmo

import (
	"fmt"
	"strconv"
	"strings"
)

func ConvertStringToValue(in string) Value {
	stringValue := strings.Trim(in, " \t")

	if i, err := strconv.ParseInt(stringValue, 0, 64); err == nil {
		return NewIntegerLiteral(i)
	}

	if f, err := strconv.ParseFloat(stringValue, 64); err == nil {
		return NewFloatLiteral(f)
	}

	return NewStringLiteral(stringValue)
}

func ConvertAnyToValue(in interface{}) Value {

	if list, canCast := in.([]interface{}); canCast {
		return ConvertListToValue(list)
	}

	if dict, canCast := in.(map[string]interface{}); canCast {
		return ConvertMapToValue(dict)
	}

	return ConvertStringToValue(fmt.Sprint(in))
}

func ConvertListToValue(in []interface{}) Value {

	list := make([]Value, len(in), len(in))

	for index, value := range in {
		list[index] = ConvertAnyToValue(value)
	}
	return NewListValue(list)
}

func ConvertListOfStringsToValue(in []string) Value {

	list := make([]Value, len(in), len(in))

	for index, value := range in {
		list[index] = ConvertStringToValue(value)
	}
	return NewListValue(list)
}

func ConvertMapToValue(in map[string]interface{}) Value {

	mapping := make(map[string]Value, len(in))

	for key, value := range in {
		mapping[key] = ConvertAnyToValue(value)
	}

	return NewDictionaryValue(nil, mapping)
}
