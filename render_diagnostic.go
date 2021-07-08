package detour

import "net/http"

type DiagnosticResult struct {
	StatusCode int
	Message    string
	Header     http.Header

	DumpNonCanonicalRequestHeaders bool // no-op
}

func (this DiagnosticResult) Render(response http.ResponseWriter, _ *http.Request) {
	http.Error(response, http.StatusText(this.StatusCode)+"\n", this.StatusCode)
}
