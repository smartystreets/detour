package httptest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"strings"
)

func NewHTTPFixture() *HTTPFixture {
	return &HTTPFixture{
		RequestMethod:  http.MethodGet,
		RequestURL:     url.URL{Path: "/"},
		RequestHeaders: make(http.Header),
		Dump:           new(bytes.Buffer),
	}
}

type HTTPFixture struct {
	RequestMethod  string
	RequestURL     url.URL
	RequestBody    string
	RequestHeaders http.Header
	RequestContext context.Context

	ResponseStatus  int
	ResponseHeaders http.Header
	ResponseBody    string

	Dump *bytes.Buffer
}

func (this *HTTPFixture) SetQueryStringParameter(key, value string) {
	query := this.RequestURL.Query()
	query.Set(key, value)
	this.RequestURL.RawQuery = query.Encode()
}

func (this *HTTPFixture) SetJSONBody(body interface{}) {
	raw, _ := json.Marshal(body)
	this.RequestBody = string(raw)
}

func (this *HTTPFixture) Serve(handler http.Handler) {
	request := this.buildRequest()
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	this.collectResponse(recorder)
}
func (this *HTTPFixture) buildRequest() *http.Request {
	request := httptest.NewRequest(
		this.RequestMethod,
		this.RequestURL.String(),
		strings.NewReader(this.RequestBody),
	)
	request.Header = this.RequestHeaders
	requestDump, _ := httputil.DumpRequest(request, true)
	fmt.Fprintf(this.Dump, "REQUEST DUMP:\n%s\n\n", formatDump(">", string(requestDump)))
	this.RequestContext = request.Context()
	return request
}
func (this *HTTPFixture) collectResponse(recorder *httptest.ResponseRecorder) {
	response := recorder.Result()
	responseDump, _ := httputil.DumpResponse(response, true)
	fmt.Fprintf(this.Dump, "RESPONSE DUMP:\n%s\n\n", formatDump("<", string(responseDump)))
	body, _ := ioutil.ReadAll(response.Body)
	this.ResponseBody = string(body)
	this.ResponseStatus = response.StatusCode
	this.ResponseHeaders = response.Header
}
func formatDump(prefix, dump string) string {
	prefix = "\n" + prefix + " "
	lines := strings.Split(strings.TrimSpace(dump), "\n")
	return prefix + strings.Join(lines, prefix)
}

func (this *HTTPFixture) ResponseBodyJSON() (actual map[string]interface{}) {
	err := json.Unmarshal([]byte(this.ResponseBody), &actual)
	if err != nil {
		log.Panicln("JSON UNMARSHAL:", err)
	}
	return actual
}
