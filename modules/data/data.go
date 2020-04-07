package data

import (
	"bytes"
	"encoding/csv"
	"fmt"

	elmo "github.com/okke/elmo/core"
)

// Module contains functions that makes the handling of structured data more easy
//
var Module = elmo.NewModule("data", initModule)

func initModule(context elmo.RunContext) elmo.Value {
	return elmo.NewMappingForModule(context, []elmo.NamedValue{
		fromCSV(), toCSV(), fromJSON(), toJSON()})
}

func fromCSV() elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("fromCSV", `converts comma separated values into a list of dictionaries
    Usage: fromCSV <string>
	Returns: list of dictionaries
	
	fromCSV will read the contents of a multiline input string. It assumes the first line
	contains the fielnames which are used to construct properties of the dictionary objects.

	example: 

	data: (load data)
	puts (((file "./test/peppers.csv") string) |data.fromCSV)

    `,
		func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

			_, err := elmo.CheckArguments(arguments, 1, 1, "fromCSV", "<string>")
			if err != nil {
				return err
			}

			value := elmo.EvalArgument(context, arguments[0])
			return convertCSVStringToListOfDictionaries(value.String())

		})
}

func toCSV() elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("toCSV", `converts an elmo list of dictionaries to csv
    Usage: toCSV <list of header strings> <list of dictionaries>
	Returns: CSV representation of value

    `,
		func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

			if _, err := elmo.CheckArguments(arguments, 2, 2, "toCSV", "<list of header strings> <list of dictionaries>"); err != nil {
				return err
			}

			headersValue := elmo.EvalArgument(context, arguments[0])
			if headersValue.Type() != elmo.TypeList {
				return elmo.NewErrorValue("toCSV: expect a list of headers")
			}

			var buf bytes.Buffer
			writer := csv.NewWriter(&buf)

			headers := make([]string, 0, 0)
			for _, header := range headersValue.(elmo.ListValue).List() {
				headers = append(headers, header.String())
			}
			writer.Write(headers)

			dataValue := elmo.EvalArgument(context, arguments[1])
			if dataValue.Type() != elmo.TypeList {
				return elmo.NewErrorValue("toCSV: expect a list with data")
			}
			data := make([]string, len(headers), len(headers))
			for _, record := range dataValue.(elmo.ListValue).List() {
				if record.Type() == elmo.TypeDictionary {
					flat := elmo.ConvertDictionaryToFlatMap(record.(elmo.DictionaryValue))

					for i, header := range headers {
						if value, found := flat[header]; found {
							data[i] = value.String()
						} else {
							data[i] = ""
						}
					}

					writer.Write(data)
				} else {
					return elmo.NewErrorValue(fmt.Sprintf("toCSV: non dictionary (%v) of type (%s) found in list of data", record, record.Info().Name()))
				}
			}

			writer.Flush()
			return elmo.NewStringLiteral(buf.String())

		})
}

func fromJSON() elmo.NamedValue {
	return elmo.NewGoFunction(`fromJSON/converts json into a dictionary
    Usage: fromJSON <string>
	Returns: Dictionary representation of given json 

	example: 

	data: (load data)
	puts (((file "./test/habanero.json") string) |data.fromJSON)

    `,
		func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

			if _, err := elmo.CheckArguments(arguments, 1, 1, "fromJSON", "<string>"); err != nil {
				return err
			}

			value := elmo.EvalArgument(context, arguments[0])
			return convertJSONStringToDictionary(value.String())

		})
}

func toJSON() elmo.NamedValue {
	return elmo.NewGoFunction(`toJSON/converts an elmo value to json
    Usage: toJSON <value>
	Returns: JSON representation of value

    `,
		func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

			if _, err := elmo.CheckArguments(arguments, 1, 1, "toJSON", "<value>"); err != nil {
				return err
			}

			return elmo.NewStringLiteral(convertValueToJSONString(elmo.EvalArgument(context, arguments[0])))

		})
}
