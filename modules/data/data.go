package data

import (
	elmo "github.com/okke/elmo/core"
)

// Module contains functions that makes the handling of structured data more easy
//
var Module = elmo.NewModule("data", initModule)

func initModule(context elmo.RunContext) elmo.Value {
	return elmo.NewMappingForModule(context, []elmo.NamedValue{
		fromCSV(), fromJSON(), toJSON()})
}

func fromCSV() elmo.NamedValue {
	return elmo.NewGoFunction(`fromCSV/converts comma separated values into a list of dictionaries
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
