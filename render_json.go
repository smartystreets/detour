package detour

import "net/http"

type JSONResult struct {
	StatusCode  int
	ContentType string
	Content     interface{}
	Indent      string
	Header      http.Header
}

func (this JSONResult) Render(response http.ResponseWriter, _ *http.Request) {
	copyHeadersToResponse(this.Header, response.Header())
	writeJSONResponse(response,
		this.StatusCode,
		this.Content,
		firstNonBlank(this.ContentType, jsonContentType),
		this.Indent)
}
