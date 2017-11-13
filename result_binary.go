package detour

import "net/http"

type BinaryResult struct {
	StatusCode  int
	ContentType string
	Content     []byte
}

func (this BinaryResult) Render(response http.ResponseWriter, request *http.Request) {
	contentType := firstNonBlank(this.ContentType, octetStreamContentType)
	writeContentTypeAndStatusCode(response, this.StatusCode, contentType)
	response.Write(this.Content)
}
