package elmohttp

import (
	elmo "github.com/okke/elmo/core"
)

// Module contains http related functions
//
var Module = elmo.NewModule("http", initModule)

func initModule(context elmo.RunContext) elmo.Value {
	return elmo.NewMappingForModule(context, []elmo.NamedValue{
		client(), get(), cookies(), testServer(), testURL()})
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

func get() elmo.NamedValue {
	return elmo.NewGoFunction(`get/executes an GET request on an http client
	`, func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {

		argLen, err := elmo.CheckArguments(arguments, 1, 2, "get", "<client> <path>")
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
