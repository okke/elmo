package elmohttp

import (
	"net/http"
	"io/ioutil"
	"fmt"
	elmo "github.com/okke/elmo/core"
)

var typeInfoHTTPClient = elmo.NewTypeInfo("httpClient")

type httpClient struct {
	client *http.Client
	baseUrl string
}

type HTTPClient interface {
	DoRequest(method string, url string) elmo.Value
}

func NewHTTPClient(baseUrl string) HTTPClient {
	return &httpClient{baseUrl: baseUrl, client: &http.Client{}}
}

func (httpClient *httpClient) String() string {
	return httpClient.baseUrl
}

func (httpClient *httpClient) DoRequest(method, url string) elmo.Value {

	req, err := http.NewRequest(method, httpClient.baseUrl + url, nil)
	if err != nil {
		return elmo.NewErrorValue(err.Error())
	}
	resp, err := httpClient.client.Do(req)
	if err != nil {
		return elmo.NewErrorValue(err.Error())
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return elmo.NewErrorValue(err.Error())
	}

	if 200 != resp.StatusCode {
		return elmo.NewErrorValue(fmt.Sprintf("http status code %d", resp.StatusCode))
	}

	return elmo.NewStringLiteral(string(body))
}
