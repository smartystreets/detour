package detour

import "net/http"

type ValidationResult struct {
	Failure1 error
	Failure2 error
	Failure3 error
	Failure4 error
}

func (this ValidationResult) Render(response http.ResponseWriter, request *http.Request) {
	writeContentType(response, jsonContentType)

	var failures Errors
	failures = failures.Append(this.Failure1)
	failures = failures.Append(this.Failure2)
	failures = failures.Append(this.Failure3)
	failures = failures.Append(this.Failure4)

	serializeAndWrite(response, http.StatusUnprocessableEntity, failures)
}
