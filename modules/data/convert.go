package data

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	elmo "github.com/okke/elmo/core"
)

func convertStringToValue(in string) elmo.Value {
	stringValue := strings.Trim(in, " \t")

	if i, err := strconv.ParseInt(stringValue, 0, 64); err == nil {
		return elmo.NewIntegerLiteral(i)
	}

	if f, err := strconv.ParseFloat(stringValue, 64); err == nil {
		return elmo.NewFloatLiteral(f)
	}

	return elmo.NewStringLiteral(stringValue)
}

func convertAnyToValue(in interface{}) elmo.Value {

	if list, canCast := in.([]interface{}); canCast {
		return convertListToValue(list)
	}

	if dict, canCast := in.(map[string]interface{}); canCast {
		return convertMapToValue(dict)
	}

	return convertStringToValue(fmt.Sprint(in))
}

func convertListToValue(in []interface{}) elmo.Value {

	list := make([]elmo.Value, len(in), len(in))

	for index, value := range in {
		list[index] = convertAnyToValue(value)
	}
	return elmo.NewListValue(list)
}

func convertMapToValue(in map[string]interface{}) elmo.Value {

	mapping := make(map[string]elmo.Value, len(in))

	for key, value := range in {
		mapping[key] = convertAnyToValue(value)
	}

	return elmo.NewDictionaryValue(nil, mapping)
}

func convertCSVStringToListOfDictionaries(in string) elmo.Value {

	r := csv.NewReader(strings.NewReader(in))

	records, err := r.ReadAll()
	if err != nil {
		return elmo.NewErrorValue(err.Error())
	}

	if len(records) == 0 {
		return elmo.NewListValue([]elmo.Value{})
	}

	list := make([]elmo.Value, 0, 0)

	// get headers and trim them
	//
	header := records[0]
	for i, h := range header {
		header[i] = strings.Trim(h, " \t")
	}

	for recordIndex, record := range records {
		if recordIndex == 0 {
			continue
		}

		mapping := make(map[string]elmo.Value, 0)
		for fieldIndex, fieldValue := range record {
			if fieldIndex >= len(header) {
				continue
			}

			fieldName := header[fieldIndex]
			if fieldName != "" {
				mapping[header[fieldIndex]] = convertStringToValue(fieldValue)
			}
		}
		list = append(list, elmo.NewDictionaryValue(nil, mapping))
	}

	return elmo.NewListValue(list)
}

func convertJSONStringToDictionary(in string) elmo.Value {

	jsonMap := make(map[string]interface{})
	if err := json.Unmarshal([]byte(in), &jsonMap); err != nil {
		return elmo.NewErrorValue(err.Error())
	}

	return convertMapToValue(jsonMap)
}
