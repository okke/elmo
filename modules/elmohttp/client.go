package elmohttp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	elmo "github.com/okke/elmo/core"
)

var typeInfoHTTPClient = elmo.NewTypeInfo("httpClient")

type httpClient struct {
	client  *http.Client
	baseUrl *url.URL
	headers map[string]string
}

type HTTPClient interface {
	DoRequest(method string, url string, body []byte) elmo.Value
	Cookies() []*http.Cookie
	SetHeaders(map[string]string)
}

func NewHTTPClient(baseUrl string) (HTTPClient, elmo.ErrorValue) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, elmo.NewErrorValue(err.Error())
	}

	url, err := url.Parse(baseUrl)
	if err != nil {
		return nil, elmo.NewErrorValue(err.Error())
	}
	return &httpClient{baseUrl: url, client: &http.Client{
		Jar: jar}}, nil
}

func (httpClient *httpClient) String() string {
	return httpClient.baseUrl.String()
}

func (httpClient *httpClient) addHeaders(req *http.Request) {
	if httpClient.headers == nil {
		return
	}

	for k, v := range httpClient.headers {
		req.Header.Set(k, v)
	}
}

func (httpClient *httpClient) DoRequest(method, url string, body []byte) elmo.Value {

	req, err := http.NewRequest(method, httpClient.baseUrl.String()+url, bytes.NewBuffer(body))
	if err != nil {
		return elmo.NewErrorValue(err.Error())
	}

	httpClient.addHeaders(req)

	resp, err := httpClient.client.Do(req)
	if err != nil {
		return elmo.NewErrorValue(err.Error())
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return elmo.NewErrorValue(err.Error())
	}

	httpClient.client.Jar.SetCookies(httpClient.baseUrl, resp.Cookies())

	if 200 != resp.StatusCode {
		return elmo.NewErrorValue(fmt.Sprintf("http status code %d", resp.StatusCode))
	}

	return elmo.NewStringLiteral(string(responseBody))
}

func (httpClient *httpClient) Cookies() []*http.Cookie {
	return httpClient.client.Jar.Cookies(httpClient.baseUrl)
}

func (httpClient *httpClient) SetHeaders(headers map[string]string) {
	httpClient.headers = headers
}
