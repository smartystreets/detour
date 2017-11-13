package detour

import "net/http"

type ContentResult struct {
	StatusCode  int
	ContentType string
	Content     string
	Headers     map[string]string // TODO: do we even need/use this?
}

func (this ContentResult) Render(response http.ResponseWriter, request *http.Request) {
	contentType := firstNonBlank(this.ContentType, plaintextContentType)

	headers := response.Header()
	for key, value := range this.Headers {
		if len(key) > 0 {
			headers.Set(key, value)
		}
	}

	writeContentTypeAndStatusCode(response, this.StatusCode, contentType)
	response.Write([]byte(this.Content))
}
