package detour

import "net/http"

type RedirectResult struct {
	Location   string
	StatusCode int
}

func (this RedirectResult) Render(response http.ResponseWriter, request *http.Request) {
	http.Redirect(response, request, this.Location, this.StatusCode)
}
