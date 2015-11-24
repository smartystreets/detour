package binding

import (
	"errors"
	"net/http"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/httptest2"
)

type ResultFixture struct {
	*gunit.Fixture

	response *httptest2.ResponseRecorder
}

func (this *ResultFixture) Setup() {
	this.response = httptest2.NewRecorder()
}

///////////////////////////////////////////////////////////////////////////////

func (this *ResultFixture) TestStatusCodeResult() {
	result := &StatusCodeResult{
		StatusCode: 456,
		Message:    "Status 456",
	}

	this.render(result)

	this.assertStatusCode(456)
	this.assertContent("Status 456")
	this.assertHasHeader("Content-Type", "text/plain")
}

func (this *ResultFixture) TestContentResult() {
	result := &ContentResult{
		StatusCode: 456,
		Content:    "Hello, World!",
	}

	this.render(result)

	this.assertStatusCode(456)
	this.assertContent("Hello, World!")
	this.assertHasHeader("Content-Type", "text/plain")
}
func (this *ResultFixture) TestContentResult_WithCustomContentType() {
	result := &ContentResult{
		StatusCode:  456,
		ContentType: "application/custom-text",
		Content:     "Hello, World!",
	}

	this.render(result)

	this.assertStatusCode(456)
	this.assertContent("Hello, World!")
	this.assertHasHeader("Content-Type", "application/custom-text")
}

func (this *ResultFixture) TestBinaryResult() {
	result := &BinaryResult{
		StatusCode: 456,
		Content:    []byte("Hello, World!"),
	}

	this.render(result)

	this.assertStatusCode(456)
	this.assertContent("Hello, World!")
	this.assertHasHeader("Content-Type", "application/octet-stream")
}
func (this *ResultFixture) TestBinaryResult_WithCustomContentType() {
	result := &BinaryResult{
		StatusCode:  456,
		ContentType: "application/custom-binary",
		Content:     []byte("Hello, World!"),
	}

	this.render(result)

	this.assertStatusCode(456)
	this.assertContent("Hello, World!")
	this.assertHasHeader("Content-Type", "application/custom-binary")
}

func (this *ResultFixture) TestJSONResult() {
	result := &JSONResult{
		StatusCode: 123,
		Content:    map[string]string{"key": "value"},
	}

	this.render(result)

	this.assertStatusCode(123)
	this.assertContent(`{"key":"value"}`)
	this.assertHasHeader("Content-Type", "application/json; charset=utf-8")
}
func (this *ResultFixture) TestJSONResult_WithCustomContentType() {
	result := &JSONResult{
		StatusCode:  123,
		ContentType: "application/custom-json",
		Content:     map[string]string{"key": "value"},
	}

	this.render(result)

	this.assertStatusCode(123)
	this.assertContent(`{"key":"value"}`)
	this.assertHasHeader("Content-Type", "application/custom-json")
}
func (this *ResultFixture) TestJSONResult_SerializationFailure_HTTP500WithErrorMessage() {
	result := &JSONResult{
		StatusCode: 123,
		Content:    new(BadJSON),
	}
	this.render(result)

	this.assertStatusCode(500)
	this.assertHasHeader("Content-Type", "application/json; charset=utf-8")
	this.assertContent(`[{"fields":["HTTP Response"],"message":"Marshal failure"}]`)
}

func (this *ResultFixture) TestValidationResult() {
	result := &ValidationResult{
		Failure1: SimpleValidationError("message1", "field1"),
		Failure2: SimpleValidationError("message2", "field2"),
		Failure3: nil,
		Failure4: ComplexValidationError("message3", "field3", "field4"),
	}

	this.render(result)

	this.assertStatusCode(422)
	this.assertContent(`[{"fields":["field1"],"message":"message1"},{"fields":["field2"],"message":"message2"},{"fields":["field3","field4"],"message":"message3"}]`)
	this.assertHasHeader("Content-Type", "application/json; charset=utf-8")
}
func (this *ResultFixture) TestValidationResult_SerializationFailure_HTTP500WithErrorMessage() {
	result := &ValidationResult{
		Failure1: new(BadJSON),
	}

	this.render(result)

	this.assertStatusCode(500)
	this.assertHasHeader("Content-Type", "application/json; charset=utf-8")
	this.assertContent(`[{"fields":["HTTP Response"],"message":"Marshal failure"}]`)
}

func (this *ResultFixture) TestCookieResult() {
	result := &CookieResult{
		Cookie1: &http.Cookie{Name: "a", Value: "1"},
		Cookie2: &http.Cookie{Name: "b", Value: "2"},
		Cookie3: nil,
		Cookie4: &http.Cookie{Name: "d", Value: "4"},
	}

	this.render(result)

	this.assertStatusCode(200)
	this.So(this.response.Header()["Set-Cookie"], should.Resemble, []string{"a=1", "b=2", "d=4"})
	this.assertContent("")
}

///////////////////////////////////////////////////////////////////////////////

func (this *ResultFixture) render(result Renderer) {
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
	this.So(this.response.HeaderMap[key], should.Resemble, []string{value})
}

///////////////////////////////////////////////////////////////////////////////

type BadJSON struct{}

func (this *BadJSON) Error() string                { return "Implement the error interface." }
func (this *BadJSON) MarshalJSON() ([]byte, error) { return nil, errors.New("GOPHERS!") }

