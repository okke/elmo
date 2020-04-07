package elmo

import (
	"errors"
	"fmt"
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

func ConvertValueToInterface(in Value, ignore ...Type) interface{} {
	return convertValueToInterface(make(map[uuid.UUID]bool, 0), in, TypeMap(ignore...))
}

func convertValueToInterface(converted map[uuid.UUID]bool, in Value, ignore map[Type]bool) interface{} {

	uuid := in.(Value).UUID()
	if _, done := converted[uuid]; done {
		return errors.New("circular data found in dictionary")
	}
	converted[uuid] = true

	if _, found := ignore[in.Type()]; found {
		return nil
	}

	switch in.Type() {
	case TypeDictionary:
		return convertDictionaryToMap(converted, in.(DictionaryValue), ignore)
	case TypeList:
		return convertListToArray(converted, in.(ListValue), ignore)
	case TypeString:
		return string(in.Internal().([]rune))
	default:
		return in.Internal()
	}
}

func ConvertListToArray(in ListValue, ignore ...Type) []interface{} {
	return convertListToArray(make(map[uuid.UUID]bool, 0), in, TypeMap(ignore...))
}

func convertListToArray(converted map[uuid.UUID]bool, in ListValue, ignore map[Type]bool) []interface{} {

	values := in.List()
	if values == nil || len(values) == 0 {
		return make([]interface{}, 0, 0)
	}

	array := make([]interface{}, len(values), len(values))
	for i, value := range values {
		array[i] = convertValueToInterface(converted, value, ignore)
	}

	return array
}

func ConvertDictionaryToMap(in DictionaryValue, ignore ...Type) map[string]interface{} {
	return convertDictionaryToMap(make(map[uuid.UUID]bool, 0), in, TypeMap(ignore...))
}

func convertDictionaryToMap(converted map[uuid.UUID]bool, in DictionaryValue, ignore map[Type]bool) map[string]interface{} {

	mapping := make(map[string]interface{})

	for _, key := range in.Keys() {

		value, _ := in.Resolve(key)
		if _, found := ignore[value.Type()]; !found {
			mapping[key] = convertValueToInterface(converted, value, ignore)
		}

	}

	return mapping
}

type keyValue struct {
	key   string
	value Value
}

func prefix(pre string, name string) string {
	if pre == "" {
		return name
	}
	return pre + "_" + name
}

func ConvertDictionaryToFlatMap(in DictionaryValue) map[string]Value {
	mapping := make(map[string]Value, 0)

	for _, tuple := range flattenDictionary("", make(map[uuid.UUID]bool, 0), in, make([]*keyValue, 0, 0)) {
		mapping[tuple.key] = tuple.value
	}

	return mapping
}

func ConvertDictionaryToFlatDictionary(in DictionaryValue) DictionaryValue {

	return NewDictionaryValue(nil, ConvertDictionaryToFlatMap(in))
}

func flattenList(pre string, visited map[uuid.UUID]bool, in ListValue, done []*keyValue) []*keyValue {
	uuid := in.(Value).UUID()
	if _, done := visited[uuid]; done {
		return []*keyValue{&keyValue{key: prefix(pre, "error"), value: NewErrorValue("can not flatten circular dictionary")}}
	}
	visited[uuid] = true

	result := done

	for i, value := range in.List() {
		if value.Type() == TypeDictionary {
			result = flattenDictionary(prefix(pre, strconv.Itoa(i)), visited, value.(DictionaryValue), result)
		} else if value.Type() == TypeList {
			result = flattenList(prefix(pre, strconv.Itoa(i)), visited, value.(ListValue), result)
		} else {
			result = append(result, &keyValue{key: prefix(pre, strconv.Itoa(i)), value: value})
		}
	}
	return result
}

func flattenDictionary(pre string, visited map[uuid.UUID]bool, in DictionaryValue, done []*keyValue) []*keyValue {
	uuid := in.(Value).UUID()
	if _, done := visited[uuid]; done {
		return []*keyValue{&keyValue{key: prefix(pre, "error"), value: NewErrorValue("can not flatten circular dictionary")}}
	}
	visited[uuid] = true

	result := done
	for _, key := range in.Keys() {

		value, _ := in.Resolve(key)
		if value.Type() == TypeDictionary {
			result = flattenDictionary(prefix(pre, key), visited, value.(DictionaryValue), result)
		} else if value.Type() == TypeList {
			result = flattenList(prefix(pre, key), visited, value.(ListValue), result)
		} else {
			result = append(result, &keyValue{key: prefix(pre, key), value: value})
		}
	}
	return result
}
