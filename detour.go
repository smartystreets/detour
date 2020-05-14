package detour

import (
	"context"
	"net/http"
)

func New(detour func() Detour, handler Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		detour := detour()
		messages := detour.Bind(request)
		if len(messages) > 0 {
			handler.Handle(request.Context(), messages...)
		}
		detour.Render(response, request)
	})
}

type Detour interface {
	Bind(*http.Request) []interface{}
	Render(http.ResponseWriter, *http.Request)
}

type Handler interface {
	Handle(ctx context.Context, messages ...interface{})
}
