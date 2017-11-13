package detour

import "net/http"

type JSONResult struct {
	StatusCode  int
	ContentType string
	Content     interface{}
}

func (this JSONResult) Render(response http.ResponseWriter, request *http.Request) {
	contentType := firstNonBlank(this.ContentType, jsonContentType)
	writeContentType(response, contentType)
	serializeAndWrite(response, this.StatusCode, this.Content)
}

type JSONPResult struct {
	StatusCode  int
	ContentType string
	Content     interface{}
}

func (this JSONPResult) Render(response http.ResponseWriter, request *http.Request) {
	contentType := firstNonBlank(this.ContentType, jsonContentType)
	writeContentType(response, contentType)
	callbackLabel := request.URL.Query().Get("callback") // We don't call request.ParseForm in every case so using the URL.Query() is safer.
	serializeAndWriteJSONP(response, this.StatusCode, this.Content, callbackLabel)
}
