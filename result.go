package detour

import (
	"encoding/json"
	"net/http"
)

type (
	StatusCodeResult struct {
		StatusCode int
		Message    string
	}
	ContentResult struct {
		StatusCode  int
		ContentType string
		Content     string
	}
	BinaryResult struct {
		StatusCode  int
		ContentType string
		Content     []byte
	}
	JSONResult struct {
		StatusCode  int
		ContentType string
		Content     interface{}
	}
	ValidationResult struct {
		Failure1 error
		Failure2 error
		Failure3 error
		Failure4 error
	}
	ErrorResult struct {
		StatusCode int
		Error1 error
		Error2 error
		Error3 error
		Error4 error
	}
	CookieResult struct {
		Cookie1 *http.Cookie
		Cookie2 *http.Cookie
		Cookie3 *http.Cookie
		Cookie4 *http.Cookie
	}
)

func (this *StatusCodeResult) Render(response http.ResponseWriter, request *http.Request) {
	writeContentTypeAndStatusCode(response, this.StatusCode, plaintextContentType)
	response.Write([]byte(this.Message))
}

func (this *ContentResult) Render(response http.ResponseWriter, request *http.Request) {
	contentType := selectContentType(this.ContentType, plaintextContentType)
	writeContentTypeAndStatusCode(response, this.StatusCode, contentType)
	response.Write([]byte(this.Content))
}

func (this *BinaryResult) Render(response http.ResponseWriter, request *http.Request) {
	contentType := selectContentType(this.ContentType, octetStreamContentType)
	writeContentTypeAndStatusCode(response, this.StatusCode, contentType)
	response.Write(this.Content)
}

func (this *JSONResult) Render(response http.ResponseWriter, request *http.Request) {
	contentType := selectContentType(this.ContentType, jsonContentType)
	writeContentType(response, contentType)
	serializeAndWrite(response, this.StatusCode, this.Content)
}

func (this *ValidationResult) Render(response http.ResponseWriter, request *http.Request) {
	writeContentType(response, jsonContentType)

	var failures Errors
	failures = failures.Append(this.Failure1)
	failures = failures.Append(this.Failure2)
	failures = failures.Append(this.Failure3)
	failures = failures.Append(this.Failure4)

	serializeAndWrite(response, 422, failures)
}

func (this *ErrorResult) Render(response http.ResponseWriter, request *http.Request) {
	writeContentType(response, jsonContentType)

	var failures Errors
	failures = failures.Append(this.Error1)
	failures = failures.Append(this.Error2)
	failures = failures.Append(this.Error3)
	failures = failures.Append(this.Error4)

	serializeAndWrite(response, this.StatusCode, failures)
}

func (this *CookieResult) Render(response http.ResponseWriter, request *http.Request) {
	for _, cookie := range []*http.Cookie{this.Cookie1, this.Cookie2, this.Cookie3, this.Cookie4} {
		if cookie != nil {
			http.SetCookie(response, cookie)
		}
	}
}

func selectContentType(values ...string) string {
	for _, value := range values {
		if len(value) > 0 {
			return value
		}
	}

	return ""
}

func writeContentTypeAndStatusCode(response http.ResponseWriter, statusCode int, contentType string) {
	writeContentType(response, contentType)
	response.WriteHeader(statusCode)
}
func writeContentType(response http.ResponseWriter, contentType string) {
	if len(contentType) > 0 {
		response.Header().Set(contentTypeHeader, contentType) // doesn't get written unless status code is written last!
	}
}

func serializeAndWrite(response http.ResponseWriter, statusCode int, content interface{}) {
	if content, err := json.Marshal(content); err == nil {
		writeContent(response, statusCode, content)
	} else {
		writeError(response)
	}
}
func writeContent(response http.ResponseWriter, statusCode int, content []byte) {
	response.WriteHeader(statusCode)
	response.Write(content)
}
func writeError(response http.ResponseWriter) {
	response.WriteHeader(http.StatusInternalServerError)
	errContent := make(Errors, 0).Append(SimpleInputError("Marshal failure", "HTTP Response"))
	content, _ := json.Marshal(errContent)
	response.Write(content)
}

const (
	contentTypeHeader      = "Content-Type"
	jsonContentType        = "application/json; charset=utf-8"
	octetStreamContentType = "application/octet-stream"
	plaintextContentType   = "text/plain"
)
