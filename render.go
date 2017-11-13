package detour

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func firstNonBlank(values ...string) string {
	for _, value := range values {
		if len(value) > 0 {
			return value
		}
	}

	return ""
}

func writeContentTypeAndStatusCode(response http.ResponseWriter, statusCode int, contentType string) {
	writeContentType(response, contentType)
	response.WriteHeader(orOK(statusCode))
}
func writeContentType(response http.ResponseWriter, contentType string) {
	if len(contentType) > 0 {
		response.Header().Set(contentTypeHeader, contentType) // doesn't get written unless status code is written last!
	}
}

func serializeAndWrite(response http.ResponseWriter, statusCode int, content interface{}) {
	if serialized, err := json.Marshal(content); err == nil {
		writeContent(response, statusCode, serialized)
	} else {
		writeInternalServerError(response)
	}
}
func serializeAndWriteJSONP(response http.ResponseWriter, statusCode int, content interface{}, label string) {
	if len(label) == 0 {
		serializeAndWrite(response, statusCode, content)
	} else if serialized, err := json.Marshal(content); err == nil {
		buffer := bytes.NewBufferString(label)
		buffer.WriteString("(")
		buffer.Write(serialized)
		buffer.WriteString(")")
		writeContent(response, statusCode, buffer.Bytes())
	} else {
		writeInternalServerError(response)
	}
}
func writeContent(response http.ResponseWriter, statusCode int, content []byte) {
	response.WriteHeader(orOK(statusCode))
	response.Write(content)
}
func writeInternalServerError(response http.ResponseWriter) {
	response.WriteHeader(http.StatusInternalServerError)
	errContent := make(Errors, 0).Append(SimpleInputError("Marshal failure", "HTTP Response"))
	content, _ := json.Marshal(errContent)
	response.Write(content)
}

func orOK(statusCode int) int {
	if statusCode == 0 {
		return http.StatusOK
	}
	return statusCode
}

const (
	contentTypeHeader      = "Content-Type"
	jsonContentType        = "application/json; charset=utf-8"
	octetStreamContentType = "application/octet-stream"
	plaintextContentType   = "text/plain; charset=utf-8"
)
