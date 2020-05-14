package detour

import (
	"context"
	"net/http"
)

func New(detour func() Detour, handler handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		detour := detour()
		handler.Handle(request.Context(), detour.Bind(request)...)
		detour.Render(response, request)
	})
}

type Detour interface {
	Bind(*http.Request) []interface{}
	Render(http.ResponseWriter, *http.Request)
}

type handler interface {
	Handle(ctx context.Context, messages ...interface{})
}
