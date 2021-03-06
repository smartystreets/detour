package detour

import "net/http"

type ErrorResult struct {
	StatusCode int
	Error1     error
	Error2     error
	Error3     error
	Error4     error
}

func (this ErrorResult) Render(response http.ResponseWriter, _ *http.Request) {
	var failures Errors
	failures = failures.Append(this.Error1)
	failures = failures.Append(this.Error2)
	failures = failures.Append(this.Error3)
	failures = failures.Append(this.Error4)

	writeJSONResponse(response, this.StatusCode, failures, jsonContentType, "")
}
