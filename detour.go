package detour

import (
	"context"
	"net/http"

	"github.com/smartystreets/detour/v3/render"
)

func New(detour func() Detour, handler handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		detour := detour()
		renderer := detour.Bind(request)
		if renderer == nil {
			handler.Handle(request.Context(), detour.MessagesToHandle()...)
			renderer = detour.Render()
		}
		if renderer != nil {
			renderer.Render(response, request)
		}
	})
}

type Detour interface {
	Bind(*http.Request) render.Renderer
	MessagesToHandle() []interface{}
	Render() render.Renderer
}

type handler interface {
	Handle(ctx context.Context, messages ...interface{})
}
