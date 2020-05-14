package detourtest

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/smartystreets/detour/v3"
)

type Fixture interface {
	Printf(string, ...interface{})
	AssertEqual(expected, actual interface{}) bool
	Assert(condition bool, messages ...string) bool
}

func Initialize(gunit Fixture) *DetourFixture {
	return &DetourFixture{
		Fixture:         gunit,
		Handler:         NewFakeHandler(),
		RequestURL:      url.URL{Path: "/"},
		RequestBody:     make(map[string]interface{}),
		RequestHeaders:  make(http.Header),
		RequestContext:  nil,
		ResponseStatus:  0,
		ResponseHeaders: nil,
		ResponseBody:    "",
	}
}

type DetourFixture struct {
	Fixture

	Handler *FakeHandler

	RequestURL     url.URL
	RequestBody    map[string]interface{}
	RequestHeaders http.Header
	RequestContext context.Context

	ResponseStatus  int
	ResponseHeaders http.Header
	ResponseBody    string
}

func (this *DetourFixture) SetQueryStringParameter(key, value string) {
	query := this.RequestURL.Query()
	query.Set(key, value)
	this.RequestURL.RawQuery = query.Encode()
}

func (this *DetourFixture) Do(callback func() detour.Detour) {
	request := this.buildRequest()
	handler := detour.New(callback, this.Handler)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	this.collectResponse(recorder)
}
func (this *DetourFixture) buildRequest() *http.Request {
	body, _ := json.Marshal(this.RequestBody)
	request := httptest.NewRequest("GET", this.RequestURL.String(), bytes.NewReader(body))
	request.Header = this.RequestHeaders
	requestDump, _ := httputil.DumpRequest(request, true)
	this.Printf("REQUEST DUMP:\n%s\n\n", formatDump(">", string(requestDump)))
	this.RequestContext = request.Context()
	return request
}
func (this *DetourFixture) collectResponse(recorder *httptest.ResponseRecorder) {
	response := recorder.Result()
	responseDump, _ := httputil.DumpResponse(response, true)
	this.Printf("RESPONSE DUMP:\n%s\n\n", formatDump("<", string(responseDump)))
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

func (this *DetourFixture) ResponseBodyJSON() (actual map[string]interface{}) {
	this.AssertEqual("application/json; charset=utf-8", this.ResponseHeaders.Get("Content-Type"))
	err := json.Unmarshal([]byte(this.ResponseBody), &actual)
	this.Assert(err == nil)
	return actual
}
