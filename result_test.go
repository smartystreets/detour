package detour

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
	this.assertHasHeader(contentTypeHeader, plaintextContentType)
}
func (this *ResultFixture) TestStatusCodeResult_StatusCodeDefaultsTo200() {
	result := &StatusCodeResult{
		StatusCode: 0,
		Message:    "Status OK",
	}

	this.render(result)

	this.assertStatusCode(200)
}

func (this *ResultFixture) TestContentResult() {
	result := &ContentResult{
		StatusCode: 456,
		Content:    "Hello, World!",
	}

	this.render(result)

	this.assertStatusCode(456)
	this.assertContent("Hello, World!")
	this.assertHasHeader(contentTypeHeader, plaintextContentType)
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
	this.assertHasHeader(contentTypeHeader, "application/custom-text")
}
func (this *ResultFixture) TestContentResult_StatusCodeDefaultsTo200() {
	result := &ContentResult{
		StatusCode: 0,
		Content:    "Status OK",
	}

	this.render(result)

	this.assertStatusCode(http.StatusOK)
}

func (this *ResultFixture) TestBinaryResult() {
	result := &BinaryResult{
		StatusCode: 456,
		Content:    []byte("Hello, World!"),
	}

	this.render(result)

	this.assertStatusCode(456)
	this.assertContent("Hello, World!")
	this.assertHasHeader(contentTypeHeader, octetStreamContentType)
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
	this.assertHasHeader(contentTypeHeader, "application/custom-binary")
}
func (this *ResultFixture) TestBinaryResult_StatusCodeDefaultsTo200() {
	result := &BinaryResult{
		StatusCode: 0,
		Content:    []byte("Status OK"),
	}

	this.render(result)

	this.assertStatusCode(http.StatusOK)
}

func (this *ResultFixture) TestJSONResult() {
	result := &JSONResult{
		StatusCode: 123,
		Content:    map[string]string{"key": "value"},
	}

	this.render(result)

	this.assertStatusCode(123)
	this.assertContent(`{"key":"value"}`)
	this.assertHasHeader(contentTypeHeader, jsonContentType)
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
	this.assertHasHeader(contentTypeHeader, "application/custom-json")
}
func (this *ResultFixture) TestJSONResult_SerializationFailure_HTTP500WithErrorMessage() {
	result := &JSONResult{
		StatusCode: 123,
		Content:    new(BadJSON),
	}
	this.render(result)

	this.assertStatusCode(500)
	this.assertHasHeader(contentTypeHeader, jsonContentType)
	this.assertContent(`[{"fields":["HTTP Response"],"message":"Marshal failure"}]`)
}
func (this *ResultFixture) TestJSONResult_StatusCodeDefaultsTo200() {
	result := &JSONResult{
		StatusCode: 0,
		Content:    42,
	}

	this.render(result)

	this.assertStatusCode(http.StatusOK)
}

func (this *ResultFixture) TestJSONPResult() {
	result := &JSONPResult{
		StatusCode: 123,
		Content:    map[string]string{"key": "value"},
		CallbackLabel: "maybe",
	}

	this.render(result)

	this.assertStatusCode(123)
	this.assertContent(`maybe({"key":"value"})`)
	this.assertHasHeader(contentTypeHeader, jsonContentType)
}
func (this *ResultFixture) TestJSONPResult_WithCustomContentType() {
	result := &JSONPResult{
		StatusCode:  123,
		ContentType: "application/custom-json",
		Content:     map[string]string{"key": "value"},
		CallbackLabel: "maybe",
	}

	this.render(result)

	this.assertStatusCode(123)
	this.assertContent(`maybe({"key":"value"})`)
	this.assertHasHeader(contentTypeHeader, "application/custom-json")
}
func (this *ResultFixture) TestJSONPResult_SerializationFailure_HTTP500WithErrorMessage() {
	result := &JSONPResult{
		StatusCode: 123,
		Content:    new(BadJSON),
		CallbackLabel: "maybe",
	}
	this.render(result)

	this.assertStatusCode(500)
	this.assertHasHeader(contentTypeHeader, jsonContentType)
	this.assertContent(`[{"fields":["HTTP Response"],"message":"Marshal failure"}]`)
}
func (this *ResultFixture) TestJSONPResult_StatusCodeDefaultsTo200() {
	result := &JSONPResult{
		StatusCode: 0,
		Content:    42,
		CallbackLabel: "maybe",
	}

	this.render(result)

	this.assertStatusCode(http.StatusOK)
}
func (this *ResultFixture) TestJSONPResult_NoCallback_SerializesAsPlainOldJSON() {
	result := &JSONPResult{
		StatusCode: 123,
		Content:    map[string]string{"key": "value"},
		CallbackLabel: "", // empty on purpose
	}

	this.render(result)

	this.assertStatusCode(123)
	this.assertContent(`{"key":"value"}`)
	this.assertHasHeader(contentTypeHeader, jsonContentType)
}


func (this *ResultFixture) TestValidationResult() {
	result := &ValidationResult{
		Failure1: SimpleInputError("message1", "field1"),
		Failure2: SimpleInputError("message2", "field2"),
		Failure3: nil,
		Failure4: CompoundInputError("message3", "field3", "field4"),
	}

	this.render(result)

	this.assertStatusCode(422)
	this.assertContent(`[{"fields":["field1"],"message":"message1"},{"fields":["field2"],"message":"message2"},{"fields":["field3","field4"],"message":"message3"}]`)
	this.assertHasHeader(contentTypeHeader, jsonContentType)
}
func (this *ResultFixture) TestValidationResult_SerializationFailure_HTTP500WithErrorMessage() {
	result := &ValidationResult{
		Failure1: new(BadJSON),
	}

	this.render(result)

	this.assertStatusCode(500)
	this.assertHasHeader(contentTypeHeader, jsonContentType)
	this.assertContent(`[{"fields":["HTTP Response"],"message":"Marshal failure"}]`)
}

func (this *ResultFixture) TestErrorResult() {
	result := &ErrorResult{
		StatusCode: 409,
		Error1:     SimpleInputError("message1", "field1"),
		Error2:     SimpleInputError("message2", "field2"),
		Error3:     nil,
		Error4:     CompoundInputError("message3", "field3", "field4"),
	}

	this.render(result)

	this.assertStatusCode(409)
	this.assertContent(`[{"fields":["field1"],"message":"message1"},{"fields":["field2"],"message":"message2"},{"fields":["field3","field4"],"message":"message3"}]`)
	this.assertHasHeader(contentTypeHeader, jsonContentType)
}
func (this *ResultFixture) TestErrorResult_StatusCodeDefaultsTo200() {
	result := &ErrorResult{
		StatusCode: 0,
		Error1:    errors.New("ok"),
	}

	this.render(result)

	this.assertStatusCode(http.StatusOK)
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
