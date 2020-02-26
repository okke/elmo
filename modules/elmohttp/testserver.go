package elmohttp

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"

	elmo "github.com/okke/elmo/core"
)

var typeInfoHTTPTestServer = elmo.NewTypeInfo("httpTestServer")

type httpTestServer struct {
	closed bool
	server *httptest.Server
}

type HTTPTestServer interface {
	Close()
	URL() string
}

// NewElmoRequestHandler converts an elmo runnable into a http request handler
//
func NewElmoRequestHandler(context elmo.RunContext, code elmo.Runnable) func(http.ResponseWriter, *http.Request) {
	return func(responseWriter http.ResponseWriter, request *http.Request) {

		requestMap := make(map[string]elmo.Value, 0)

		// copy request parameters into request value
		//
		for k, v := range request.URL.Query() {
			if len(v) == 1 {
				requestMap[k] = elmo.ConvertStringToValue(v[0])
			} else {
				requestMap[k] = elmo.ConvertListOfStringsToValue(v)
			}

		}
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		}
		if body != nil && len(body) > 0 {
			requestMap["body"] = elmo.NewStringLiteral(string(body))
		}

		requestMap["method"] = elmo.NewStringLiteral(request.Method)

		requestValue := elmo.NewDictionaryValue(nil, requestMap)

		responseValue := elmo.NewDictionaryValue(nil, map[string]elmo.Value{
			"write":      responseWrite(responseWriter),
			"sendStatus": responseWriteStatus(responseWriter),
			"sendCookie": responseWriteCookie(responseWriter)})

		arguments := []elmo.Argument{elmo.NewDynamicArgument(requestValue), elmo.NewDynamicArgument(responseValue)}

		result := code.Run(context, arguments)

		if result != nil && result.Type() == elmo.TypeError {
			http.Error(responseWriter, result.String(), http.StatusInternalServerError)
		}

	}
}

func responseWrite(responseWriter http.ResponseWriter) elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("write", "writes content to http response", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {
		for _, arg := range arguments {
			writeValue := elmo.EvalArgument(context, arg)
			var err error = nil
			if writeValue.Type() == elmo.TypeBinary {
				_, err = responseWriter.Write(writeValue.(elmo.BinaryValue).AsBytes())
			} else {
				_, err = responseWriter.Write([]byte(writeValue.String()))
			}
			if err != nil {
				return elmo.NewErrorValue(err.Error())
			}
		}
		return elmo.Nothing
	})
}

func responseWriteStatus(responseWriter http.ResponseWriter) elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("sendStatus", "writes http status code to http response", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {
		if _, err := elmo.CheckArguments(arguments, 1, 1, "sendStatus", "<status>"); err != nil {
			return err
		}

		status := elmo.EvalArgument(context, arguments[0])
		if status.Type() != elmo.TypeInteger {
			return elmo.NewErrorValue("http status must be an integer")
		}

		responseWriter.WriteHeader(int(status.Internal().(int64)))

		return elmo.Nothing
	})
}

func responseWriteCookie(responseWriter http.ResponseWriter) elmo.NamedValue {
	return elmo.NewGoFunctionWithHelp("sendCookie", "writes http status code to http response", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {
		if _, err := elmo.CheckArguments(arguments, 3, 3, "sendCookie", "<key> <value> <duration in seconds>"); err != nil {
			return err
		}

		name := elmo.EvalArgument2String(context, arguments[0])
		value := elmo.EvalArgument2String(context, arguments[1])
		duration := elmo.EvalArgument(context, arguments[2])

		if duration.Type() != elmo.TypeInteger {
			return elmo.NewErrorValue("cookie duration must be an integer")
		}

		expire := time.Now().AddDate(0, 0, int(duration.Internal().(int64)))
		cookie := http.Cookie{
			Name:    name,
			Value:   value,
			Expires: expire,
		}
		http.SetCookie(responseWriter, &cookie)

		return elmo.Nothing
	})
}

// NewTestServer constructs a
func NewTestServer(context elmo.RunContext, code elmo.Runnable) HTTPTestServer {
	return &httpTestServer{server: httptest.NewServer(http.HandlerFunc(NewElmoRequestHandler(context, code)))}
}

func (httpTestServer *httpTestServer) Close() {
	httpTestServer.server.Close()
	httpTestServer.closed = true
}

func (httpTestServer *httpTestServer) URL() string {
	if httpTestServer.closed {
		return ""
	}
	return httpTestServer.server.URL
}
