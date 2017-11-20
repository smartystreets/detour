package detour

import "net/http"

type JSONResult struct {
	StatusCode  int
	ContentType string
	Content     interface{}
	Indent      string
}

func (this JSONResult) Render(response http.ResponseWriter, request *http.Request) {
	writeJSONResponse(response,
		this.StatusCode,
		this.Content,
		firstNonBlank(this.ContentType, jsonContentType),
		this.Indent)
}
