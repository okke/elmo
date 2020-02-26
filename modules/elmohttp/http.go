package elmohttp

import (
	elmo "github.com/okke/elmo/core"
)

// Module contains http related functions
//
var Module = elmo.NewModule("http", initModule)

func initModule(context elmo.RunContext) elmo.Value {
	return elmo.NewMappingForModule(context, []elmo.NamedValue{
		client(), get(), post(), cookies(), testServer(), testURL()})
}

func client() elmo.NamedValue {
	return elmo.NewGoFunction(`client/create a new http client
	`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		_, err := elmo.CheckArguments(arguments, 1, 1, "client", "url")
		if err != nil {
			return err
		}

		url := elmo.EvalArgument2String(context, arguments[0])
		client, err := NewHTTPClient(url)
		if err != nil {
			return err
		}

		return elmo.NewInternalValue(typeInfoHTTPClient, client)
	})
}

func getPath(context elmo.RunContext, pathArg elmo.Argument, parametersArg elmo.Argument) (string, elmo.ErrorValue) {
	path := ""

	if pathArg != nil {
		path = elmo.EvalArgument2String(context, pathArg)
	}

	if parametersArg != nil {
		parameters := elmo.EvalArgument(context, parametersArg)
		if parameters.Type() != elmo.TypeDictionary {
			return "", elmo.NewErrorValue("expect a dictionary with get parameters")
		}
		path = addParametersToPath(path, parameters.(elmo.DictionaryValue))
	}

	return path, nil
}

func getPathAndParamatersArg(arguments []elmo.Argument, pathArgNo int, parametersArgNo int) (elmo.Argument, elmo.Argument) {
	var pathArg elmo.Argument = nil
	var parametersArg elmo.Argument = nil

	if len(arguments) == pathArgNo+1 {
		pathArg = arguments[pathArgNo]
	}

	if len(arguments) == parametersArgNo+1 {
		parametersArg = arguments[parametersArgNo]
	}

	return pathArg, parametersArg
}

func get() elmo.NamedValue {
	return elmo.NewGoFunction(`get/executes an GET request on an http client
	`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		_, err := elmo.CheckArguments(arguments, 1, 3, "get", "<client> <path>? <parameters>?")
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
		pathArg, parametersArg := getPathAndParamatersArg(arguments, 1, 2)

		path, err = getPath(context, pathArg, parametersArg)
		if err != nil {
			return err
		}

		return client.Internal().(HTTPClient).DoRequest("GET", path, nil)
	})
}

func post() elmo.NamedValue {
	return elmo.NewGoFunction(`post/executes an POST request on an http client
	`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		_, err := elmo.CheckArguments(arguments, 2, 4, "post", "<client> <body> <path>? <parameters>?")
		if err != nil {
			return err
		}

		// first argument is the http client
		//
		client := elmo.EvalArgument(context, arguments[0])

		if !client.IsType(typeInfoHTTPClient) {
			return elmo.NewErrorValue("invalid call to http.get, expected an http client as first parameter")
		}

		// second argument is body to post
		//
		body := []byte(elmo.EvalArgument2String(context, arguments[1]))

		path := ""
		pathArg, parametersArg := getPathAndParamatersArg(arguments, 2, 3)

		path, err = getPath(context, pathArg, parametersArg)
		if err != nil {
			return err
		}

		return client.Internal().(HTTPClient).DoRequest("POST", path, body)
	})
}

func cookies() elmo.NamedValue {
	return elmo.NewGoFunction(`cookies/return all cookies of http client
	`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {
		_, err := elmo.CheckArguments(arguments, 1, 1, "cookies", "<client>")
		if err != nil {
			return err
		}

		// first argument is the http client
		//
		client := elmo.EvalArgument(context, arguments[0])

		if !client.IsType(typeInfoHTTPClient) {
			return elmo.NewErrorValue("invalid call to http.get, expected an http client as first parameter")
		}

		mapping := make(map[string]elmo.Value, 0)
		for _, cookie := range client.Internal().(HTTPClient).Cookies() {
			mapping[cookie.Name] = elmo.NewStringLiteral(cookie.Value)
		}

		return elmo.NewDictionaryValue(nil, mapping)
	})
}

func testServer() elmo.NamedValue {
	return elmo.NewGoFunction(`testServer/create a new http test server
	`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {
		if _, err := elmo.CheckArguments(arguments, 1, 1, "testServer", "<block>"); err != nil {
			return err
		}

		handler := elmo.EvalArgument(context, arguments[0])
		if handler.Type() != elmo.TypeGoFunction {
			return elmo.NewErrorValue("test_server expects function as argument")
		}

		server := NewTestServer(context.CreateSubContext(), handler.(elmo.Runnable))

		return elmo.NewInternalValue(typeInfoHTTPTestServer, server)
	})
}

func testURL() elmo.NamedValue {
	return elmo.NewGoFunction(`testURL/retrieves the url of a given test server
	`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {
		if _, err := elmo.CheckArguments(arguments, 1, 1, "testURL", "<test_server>"); err != nil {
			return err
		}

		server := elmo.EvalArgument(context, arguments[0])
		if !server.IsType(typeInfoHTTPTestServer) {
			return elmo.NewErrorValue("test_url expects test server as argument")
		}

		return elmo.NewStringLiteral(server.Internal().(HTTPTestServer).URL())
	})
}
