package data

import (
	elmo "github.com/okke/elmo/core"
)

// Module contains functions that makes the handling of structured data more easy
//
var Module = elmo.NewModule("data", initModule)

func initModule(context elmo.RunContext) elmo.Value {
	return elmo.NewMappingForModule(context, []elmo.NamedValue{
		_csv(), _json()})
}

func _csv() elmo.NamedValue {
	return elmo.NewGoFunction(`csv/converts comma separated values into a list of dictionaries
    Usage: csv <string>
	Returns: list of dictionaries
	
	cvs will read the contents of a multiline input string. It assumes the first line
	contains the fielnames which are used to construct properties of the dictionary objects.

	example: 

	data: (load data)
	puts (((file "./test/peppers.csv") string) |data.csv)

    `,
		func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

			_, err := elmo.CheckArguments(arguments, 1, 1, "csv", "<string>")
			if err != nil {
				return err
			}

			value := elmo.EvalArgument(context, arguments[0])
			return convertCSVStringToListOfDictionaries(value.String())

		})
}

func _json() elmo.NamedValue {
	return elmo.NewGoFunction(`json/converts json into a dictionary
    Usage: json <string>
	Returns: Dictionary representation of given json (a fulll elmo Value tree)

	example: 

	data: (load data)
	puts (((file "./test/habanero.json") string) |data.json)

    `,
		func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

			if _, err := elmo.CheckArguments(arguments, 1, 1, "json", "<string>"); err != nil {
				return err
			}

			value := elmo.EvalArgument(context, arguments[0])
			return convertJSONStringToDictionary(value.String())

		})
}
