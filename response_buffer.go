package detour

import (
	"bytes"
	"io"
	"net/http"
)

type responseBuffer struct {
	statusCode int
	headers    http.Header
	body       *bytes.Buffer
}

func newResponseBuffer() *responseBuffer {
	buffer := new(responseBuffer)
	buffer.initialize()
	return buffer
}

func (this *responseBuffer) StatusCode() int             { return this.statusCode }
func (this *responseBuffer) Header() http.Header         { return this.headers }
func (this *responseBuffer) Write(p []byte) (int, error) { return this.body.Write(p) }
func (this *responseBuffer) WriteHeader(statusCode int)  { this.statusCode = statusCode }

func (this *responseBuffer) flush(response http.ResponseWriter) {
	copyHeaders(this.headers, response.Header())
	response.WriteHeader(this.statusCode)
	_, _ = io.Copy(response, this.body)
	this.initialize()
}
func copyHeaders(source, destination http.Header) {
	for key, value := range source {
		destination[key] = append(destination[key], value...)
	}
}
func (this *responseBuffer) initialize() {
	this.initializeStatusCode()
	this.initializeHeaders()
	this.initializeBody()
}
func (this *responseBuffer) initializeStatusCode() {
	this.statusCode = http.StatusOK
}
func (this *responseBuffer) initializeBody() {
	if this.body == nil {
		this.body = new(bytes.Buffer)
	} else {
		this.body.Reset()
	}
}
func (this *responseBuffer) initializeHeaders() {
	if this.headers == nil {
		this.headers = make(http.Header)
	} else {
		this.resetHeaders()
	}
}
func (this *responseBuffer) resetHeaders() {
	for key := range this.headers {
		delete(this.headers, key)
	}
}

