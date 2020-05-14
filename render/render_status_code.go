package render

import "net/http"

type StatusCodeResult struct {
	StatusCode int
	Message    string
}

func (this StatusCodeResult) Render(response http.ResponseWriter, _ *http.Request) {
	writeContentTypeAndStatusCode(response, this.StatusCode, plaintextContentType)
	response.Write([]byte(this.Message))
}
