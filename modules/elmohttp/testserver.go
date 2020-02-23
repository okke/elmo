package elmohttp

import (
	"net/http"
	"net/http/httptest"

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

		requestValue := elmo.NewDictionaryValue(nil, requestMap)

		responseValue := elmo.NewDictionaryValue(nil, map[string]elmo.Value{
			"write": elmo.NewGoFunctionWithHelp("write", "writes content to http response", func(context elmo.RunContext, arguments []elmo.Argument) elmo.Value {
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
			})})

		arguments := []elmo.Argument{elmo.NewDynamicArgument(requestValue), elmo.NewDynamicArgument(responseValue)}

		result := code.Run(context, arguments)

		if result != nil && result.Type() == elmo.TypeError {
			http.Error(responseWriter, result.String(), http.StatusInternalServerError)
		}

	}
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
