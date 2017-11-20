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

func writeJSONResponse(response http.ResponseWriter, statusCode int, content interface{}, contentType, indent string) {
	writeContentType(response, contentType)
	serialized, err := serializeJSON(content, indent)
	writeResponse(response, statusCode, serialized, err)
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

func serializeJSON(content interface{}, indent string) ([]byte, error) {
	writer := new(bytes.Buffer)
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", indent)
	if err := encoder.Encode(content); err != nil {
		return nil, err
	} else {
		return writer.Bytes(), nil
	}
}

func writeResponse(response http.ResponseWriter, statusCode int, content []byte, previous error) {
	if previous == nil {
		writeContent(response, statusCode, content)
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
	content, _ := serializeJSON(errContent, "")
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
