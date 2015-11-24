package binding

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"bytes"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type ResultFixture struct {
	*gunit.Fixture

	response *httptest.ResponseRecorder
}

func (this *ResultFixture) Setup() {
	this.response = httptest.NewRecorder()
}

///////////////////////////////////////////////////////////////////////////////

func (this *ResultFixture) TestUnmodifiedResult_HTTP200_BlankBody_NoHeaders() {
	this.render(DefaultResult())

	this.assertStatusCode(200)
	this.assertNoHeaders()
	this.assertContent("")
}

func (this *ResultFixture) TestStatusCodeResult() {
	this.render(StatusCodeResult(201))

	this.assertStatusCode(201)
	this.assertNoHeaders()
	this.assertContent("")
}

func (this *ResultFixture) TestInvalidResult() {
	this.render(InvalidResult("input-field", "It's just all wrong."))

	this.assertStatusCode(422)
	this.assertHasHeader(contentTypeHeader, jsonContentType)
	this.assertContent(`[{"fields":["input-field"],"message":"It's just all wrong."}]`)
}

func (this *ResultFixture) TestNotFoundResult() {
	this.render(NotFoundResult())

	this.assertStatusCode(http.StatusNotFound)
	this.assertHasHeader(contentTypeHeader, plainTextContentType)
	this.assertContent(http.StatusText(http.StatusNotFound))
}

func (this *ResultFixture) TestStringContentResult() {
	this.render(StringContentResult("text/html", "<html></html>"))

	this.assertStatusCode(200)
	this.assertHasHeader(contentTypeHeader, "text/html")
	this.assertContent("<html></html>")
}

func (this *ResultFixture) TestContentResult() {
	this.render(ContentResult("text/html", []byte("<html></html>")))

	this.assertStatusCode(200)
	this.assertHasHeader(contentTypeHeader, "text/html")
	this.assertContent("<html></html>")
}

func (this *ResultFixture) TestHeaderResult() {
	this.render(HeaderResult("Key", "Value"))

	this.assertStatusCode(200)
	this.assertHasHeader("Key", "Value")
	this.assertContent("")
}

func (this *ResultFixture) TestCookieResult() {
	this.render(CookieResult(&http.Cookie{
		Domain: "domain.com",
		Name:   "cookie-name",
		Path:   "/path",
	}))

	this.assertStatusCode(200)
	this.assertHasHeader("Set-Cookie", "cookie-name=; Path=/path; Domain=domain.com")
	this.assertContent("")
}

func (this *ResultFixture) TestJSONResult() {
	this.render(JSONResult(struct {
		Message string `json:"message"`
	}{Message: "Hello, World!"}))

	this.assertStatusCode(200)
	this.assertHasHeader(contentTypeHeader, jsonContentType)
	this.assertContent(`{"message":"Hello, World!"}`)
}

func (this *ResultFixture) TestJSONSerializationError_HTTP500() {
	response := NewErrorProneResponseWriter()
	JSONResult(struct{}{}).Render(response, nil)
	this.So(response.StatusCode, should.Equal, 500)
	this.So(response.buffer.String(), should.Equal, "Response serialization failed")
	this.So(response.headers.Get("Content-Type"), should.Equal, plainTextContentType)
}

func (this *ResultFixture) TestErrorsPreventJSONFromBeingSerializedAndReturned() {
	this.render(DefaultResult().
		SetJSONContent(struct{ Message string }{"Not a chance"}).
		AppendInvalidResult("blah", "blah"))
	this.assertStatusCode(422)
	this.assertHasHeader(contentTypeHeader, jsonContentType)
	this.assertContent(`[{"fields":["blah"],"message":"blah"}]`)
}

///////////////////////////////////////////////////////////////////////////////

func (this *ResultFixture) render(result *Result) {
	result.Render(this.response, nil)
}
func (this *ResultFixture) assertStatusCode(expected int) {
	this.So(this.response.Code, should.Equal, expected)
}
func (this *ResultFixture) assertContent(expected string) {
	this.So(this.response.Body.String(), should.EqualTrimSpace, expected)
}
func (this *ResultFixture) assertNoHeaders() {
	this.So(this.response.Header(), should.HaveLength, 0)
}
func (this *ResultFixture) assertHasHeader(key, value string) {
	this.So(this.response.HeaderMap, should.ContainKey, key)
	this.So(this.response.Header().Get(key), should.Equal, value)
}

///////////////////////////////////////////////////////////////////////////////

type ErrorProneResponseRecorder struct {
	StatusCode int
	calls      int
	buffer     *bytes.Buffer
	headers    http.Header
}

func NewErrorProneResponseWriter() *ErrorProneResponseRecorder {
	return &ErrorProneResponseRecorder{
		buffer:  new(bytes.Buffer),
		headers: make(http.Header),
	}
}

func (this *ErrorProneResponseRecorder) Write(p []byte) (int, error) {
	this.calls++
	if this.calls == 1 {
		return 0, errors.New("GOPHERS!")
	}
	return this.buffer.Write(p)
}

func (this *ErrorProneResponseRecorder) WriteHeader(code int) {
	this.StatusCode = code
}

func (this *ErrorProneResponseRecorder) Header() http.Header {
	return this.headers
}

///////////////////////////////////////////////////////////////////////////////
