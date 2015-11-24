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
	writeStatusAndContentType(response, this.StatusCode, "text/plain")
	if len(this.Message) > 0 {
		response.Write([]byte(this.Message))
	}
}

func (this *ContentResult) Render(response http.ResponseWriter, request *http.Request) {
	writeStatusAndContentType(response, this.StatusCode, selectContentType(this.ContentType, "text/plain"))
	response.Write([]byte(this.Content))
}

func (this *BinaryResult) Render(response http.ResponseWriter, request *http.Request) {
	writeStatusAndContentType(response, this.StatusCode, selectContentType(this.ContentType, "application/octet-stream"))
	response.Write(this.Content)
}

func (this *JSONResult) Render(response http.ResponseWriter, request *http.Request) {
	writeStatusAndContentType(response, this.StatusCode, selectContentType(this.ContentType, "application/json; charset=utf-8"))
	json.NewEncoder(response).Encode(this.Content)
}

func (this *ValidationResult) Render(response http.ResponseWriter, request *http.Request) {
	writeStatusAndContentType(response, 422, "application/json; charset=utf-8")

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
	response.WriteHeader(statusCode)
	if len(contentType) > 0 {
		response.Header().Set("Content-Type", contentType)
	}
}
