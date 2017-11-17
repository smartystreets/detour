package detour

import "net/http"

type JSONResult struct {
	StatusCode  int
	ContentType string
	Content     interface{}
	Indent      string
}

func (this JSONResult) Render(response http.ResponseWriter, request *http.Request) {
	writeContentType(response, firstNonBlank(this.ContentType, jsonContentType))
	content, err := serializeJSON(this.Content, this.Indent)
	writeResponse(response, this.StatusCode, content, err)
}
