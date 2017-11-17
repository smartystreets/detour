package detour

import "net/http"

type ErrorResult struct {
	StatusCode int
	Error1     error
	Error2     error
	Error3     error
	Error4     error
}

func (this ErrorResult) Render(response http.ResponseWriter, request *http.Request) {
	writeContentType(response, jsonContentType)

	var failures Errors
	failures = failures.Append(this.Error1)
	failures = failures.Append(this.Error2)
	failures = failures.Append(this.Error3)
	failures = failures.Append(this.Error4)

	content, err := serializeJSON(failures, "")
	writeResponse(response, this.StatusCode, content, err)
}
