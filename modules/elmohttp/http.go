package elmohttp

import "github.com/okke/elmo/core"

// Module contains http related functions 
//
var Module = elmo.NewModule("http", initModule)

func initModule(context elmo.RunContext) elmo.Value {
	return elmo.NewMappingForModule(context, []elmo.NamedValue{
		client(), get()})
}

func client() elmo.NamedValue {
	return elmo.NewGoFunction(`client/create a new http client
	`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		_, err := elmo.CheckArguments(arguments, 1, 1, "client", "url")
		if err != nil {
			return err
		}

		url := elmo.EvalArgument2String(context, arguments[0])
		client := elmo.NewInternalValue(typeInfoHTTPClient, NewHTTPClient(url))

		return client
	})
}

func get() elmo.NamedValue {
	return elmo.NewGoFunction(`get/executes an GET request on an http client
	`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		argLen, err := elmo.CheckArguments(arguments, 1, 2, "get", "client <path>")
		if err != nil {
			return err
		}

		// first argument is the http client
		//
		client := elmo.EvalArgument(context, arguments[0])

		if !client.IsType(typeInfoHTTPClient) {
			return elmo.NewErrorValue("invalid call to http.get, expected an http client as first parameter")
		}

		path := ""
		if argLen == 2 {
			path = elmo.EvalArgument2String(context, arguments[1])
		}

		return client.Internal().(HTTPClient).DoRequest("GET", path)
	})
}



