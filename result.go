package binding

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
	CookieResult struct {
		Cookie1 *http.Cookie
		Cookie2 *http.Cookie
		Cookie3 *http.Cookie
		Cookie4 *http.Cookie
	}
)

func (this *StatusCodeResult) Render(response http.ResponseWriter, request *http.Request) {
	writeStatusAndContentType(response, this.StatusCode, plaintextContentType)
	if len(this.Message) > 0 {
		response.Write([]byte(this.Message))
	}
}

func (this *ContentResult) Render(response http.ResponseWriter, request *http.Request) {
	contentType := selectContentType(this.ContentType, plaintextContentType)
	writeStatusAndContentType(response, this.StatusCode, contentType)
	response.Write([]byte(this.Content))
}

func (this *BinaryResult) Render(response http.ResponseWriter, request *http.Request) {
	contentType := selectContentType(this.ContentType, plaintextContentType)
	writeStatusAndContentType(response, this.StatusCode, contentType)
	response.Write(this.Content)
}

func (this *JSONResult) Render(response http.ResponseWriter, request *http.Request) {
	contentType := selectContentType(this.ContentType, jsonContentType)
	writeStatusAndContentType(response, this.StatusCode, contentType)
	json.NewEncoder(response).Encode(this.Content)
}

func (this *ValidationResult) Render(response http.ResponseWriter, request *http.Request) {
	writeStatusAndContentType(response, 422, jsonContentType)

	var failures ValidationErrors
	failures = failures.Append(this.Failure1)
	failures = failures.Append(this.Failure2)
	failures = failures.Append(this.Failure3)
	failures = failures.Append(this.Failure4)
	json.NewEncoder(response).Encode(failures)
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

func writeStatusAndContentType(response http.ResponseWriter, statusCode int, contentType string) {
	if len(contentType) > 0 {
		response.Header().Set(contentTypeHeader, contentType) // doesn't get written unless status code is written last!
	}
	response.WriteHeader(statusCode)
}

const (
	contentTypeHeader      = "Content-Type"
	jsonContentType        = "application/json; charset=utf-8"
	octetStreamContentType = "application/octet-stream"
	plaintextContentType   = "text/plain"
)
