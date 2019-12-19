package detour

import "net/http"

type ValidationResult struct {
	Failure1 error
	Failure2 error
	Failure3 error
	Failure4 error
}

func (this ValidationResult) Render(response http.ResponseWriter, _ *http.Request) {
	var failures Errors
	failures = failures.Append(this.Failure1)
	failures = failures.Append(this.Failure2)
	failures = failures.Append(this.Failure3)
	failures = failures.Append(this.Failure4)

	writeJSONResponse(response, http.StatusUnprocessableEntity, failures, jsonContentType, "")
}
