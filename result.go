package binding

import (
	"encoding/json"
	"net/http"
)

type Result struct {
	StatusCode int
	Content    []byte
	Invalid    ValidationErrors
	JSON       interface{}
	Headers    map[string]string
	Cookies    []*http.Cookie
}

///////////////////////////////////////////////////////////////////////////////

func DefaultResult() *Result {
	return &Result{}
}
func StatusCodeResult(code int) *Result {
	return DefaultResult().SetStatusCode(code)
}
func InvalidResult(field, message string) *Result {
	return DefaultResult().AppendInvalidResult(field, message)
}
func NotFoundResult() *Result {
	return DefaultResult().
		SetStatusCode(http.StatusNotFound).
		SetStringContent(plainTextContentType, http.StatusText(http.StatusNotFound))
}
func ContentResult(contentType string, content []byte) *Result {
	return DefaultResult().SetContent(contentType, content)
}
func StringContentResult(contentType, content string) *Result {
	return DefaultResult().SetStringContent(contentType, content)
}
func HeaderResult(key, value string) *Result {
	return DefaultResult().SetHeader(key, value)
}
func CookieResult(cookie *http.Cookie) *Result {
	return DefaultResult().SetCookie(cookie)
}
func JSONResult(content interface{}) *Result {
	return DefaultResult().SetJSONContent(content)
}

///////////////////////////////////////////////////////////////////////////////

func (this *Result) SetHeader(key, value string) *Result {
	if this.Headers == nil {
		this.Headers = make(map[string]string)
	}
	this.Headers[key] = value
	return this
}
func (this *Result) SetStringContent(contentType, content string) *Result {
	return this.SetContent(contentType, []byte(content))
}
func (this *Result) SetContent(contentType string, content []byte) *Result {
	this.Content = content
	return this.SetContentType(contentType)
}
func (this *Result) SetJSONContent(content interface{}) *Result {
	this.JSON = content
	return this
}
func (this *Result) SetStatusCode(code int) *Result {
	this.StatusCode = code
	return this
}
func (this *Result) SetContentType(value string) *Result {
	return this.SetHeader(contentTypeHeader, value)
}
func (this *Result) SetCookie(cookie *http.Cookie) *Result {
	this.Cookies = append(this.Cookies, cookie)
	return this
}
func (this *Result) AppendInvalidResult(field, message string) *Result {
	this.Invalid = this.Invalid.Append(SimpleValidationError(message, field))
	return this.SetStatusCode(422)
}

///////////////////////////////////////////////////////////////////////////////

func (this *Result) Render(response http.ResponseWriter, request *http.Request) {
	this.writeStatusCode(response)
	headers := response.Header()
	this.writeHeaders(headers, response)
	this.writeBody(headers, response)
}
func (this *Result) writeStatusCode(response http.ResponseWriter) {
	if this.StatusCode == 0 {
		this.StatusCode = http.StatusOK
	}
	response.WriteHeader(this.StatusCode)
}
func (this *Result) writeHeaders(headers http.Header, response http.ResponseWriter) {
	for key, value := range this.Headers {
		headers.Set(key, value)
	}

	for _, cookie := range this.Cookies {
		http.SetCookie(response, cookie)
	}
}
func (this *Result) writeBody(headers http.Header, response http.ResponseWriter) {
	if err := this.tryWriteBody(headers, response); err != nil {
		this.writeFailure(headers, response)
	}
}
func (this *Result) tryWriteBody(headers http.Header, response http.ResponseWriter) error {
	if len(this.Invalid) > 0 {
		return this.serializeJSON(headers, response, this.Invalid)
	} else if this.JSON != nil {
		return this.serializeJSON(headers, response, this.JSON)
	} else if len(this.Content) > 0 {
		return this.writePlainText(headers, response)
	} else {
		return nil
	}
}
func (this *Result) serializeJSON(headers http.Header, response http.ResponseWriter, content interface{}) error {
	headers.Set(contentTypeHeader, jsonContentType)
	return json.NewEncoder(response).Encode(content)
}
func (this *Result) writePlainText(headers http.Header, response http.ResponseWriter) error {
	if len(headers.Get(contentTypeHeader)) == 0 {
		headers.Set(contentTypeHeader, plainTextContentType)
	}
	_, err := response.Write(this.Content)
	return err
}
func (this *Result) writeFailure(headers http.Header, response http.ResponseWriter) {
	// TODO: log the failure
	headers.Del(contentTypeHeader)
	response.WriteHeader(http.StatusInternalServerError)
	this.Content = []byte("Response serialization failed")
	this.writePlainText(headers, response)
}

///////////////////////////////////////////////////////////////////////////////

const (
	contentTypeHeader    = "Content-Type"
	plainTextContentType = "text/plain; charset=utf-8"
	jsonContentType      = "application/json; charset=utf-8"
)
