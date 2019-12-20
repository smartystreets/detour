package detour

import (
	"fmt"
	"net/http"
	"strings"
)

type JSONPResult struct {
	StatusCode  int
	ContentType string
	Content     interface{}
	Indent      string
	Header      http.Header
}

func (this JSONPResult) Render(response http.ResponseWriter, request *http.Request) {
	copyHeaders(this.Header, response.Header())
	writeContentType(response, firstNonBlank(this.ContentType, jsonContentType))
	content, err := serializeJSON(this.Content, this.Indent)
	content = wrapJSONP(content, callbackLabel(request))
	writeResponse(response, this.StatusCode, content, err)
}

func wrapJSONP(content []byte, label string) []byte {
	serialized := strings.TrimSpace(string(content))
	if len(label) > 0 {
		serialized = fmt.Sprintf("%s(%s)", label, serialized)
	}
	return []byte(serialized)
}

func callbackLabel(request *http.Request) string {
	// We don't call request.ParseForm in every case so using the URL.Query() is safer.
	return request.URL.Query().Get("callback")
}
