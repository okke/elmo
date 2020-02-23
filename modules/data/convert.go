package data

import (
	"encoding/csv"
	"encoding/json"
	"strings"

	elmo "github.com/okke/elmo/core"
)

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
				mapping[header[fieldIndex]] = elmo.ConvertStringToValue(fieldValue)
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

	return elmo.ConvertMapToValue(jsonMap)
}
