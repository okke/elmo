package elmo

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/google/uuid"
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

func ConvertValueToInterface(in Value) interface{} {
	return convertValueToInterface(make(map[uuid.UUID]bool, 0), in)
}

func convertValueToInterface(converted map[uuid.UUID]bool, in Value) interface{} {

	uuid := in.(Value).UUID()
	if _, done := converted[uuid]; done {
		return errors.New("circular data found in dictionary")
	}
	converted[uuid] = true

	switch in.Type() {
	case TypeDictionary:
		return convertDictionaryToMap(converted, in.(DictionaryValue))
	case TypeList:
		return convertListToArray(converted, in.(ListValue))
	case TypeString:
		return string(in.Internal().([]rune))
	default:
		return in.Internal()
	}
}

func ConvertListToArray(in ListValue) []interface{} {
	return convertListToArray(make(map[uuid.UUID]bool, 0), in)
}

func convertListToArray(converted map[uuid.UUID]bool, in ListValue) []interface{} {

	values := in.List()
	if values == nil || len(values) == 0 {
		return make([]interface{}, 0, 0)
	}

	array := make([]interface{}, len(values), len(values))
	for i, value := range values {
		array[i] = convertValueToInterface(converted, value)
	}

	return array
}

func ConvertDictionaryToMap(in DictionaryValue) map[string]interface{} {
	return convertDictionaryToMap(make(map[uuid.UUID]bool, 0), in)
}

func convertDictionaryToMap(converted map[uuid.UUID]bool, in DictionaryValue) map[string]interface{} {

	mapping := make(map[string]interface{})
	keys := in.Keys()
	sort.Strings(keys)

	for _, key := range keys {

		value, _ := in.Resolve(key)
		mapping[key] = convertValueToInterface(converted, value)
	}

	return mapping
}
