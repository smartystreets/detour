package detour

import (
	"context"
	"net/http"

	"github.com/smartystreets/detour/v3/render"
)

func New(detour func() Detour, handler Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		detour := detour()
		renderer := detour.Bind(request)
		if renderer == nil {
			renderer = detour.Handle(request.Context(), handler)
		}
		if renderer != nil {
			renderer.Render(response, request)
		}
	})
}

type Detour interface {
	Bind(*http.Request) render.Renderer
	Handle(context.Context, Handler) render.Renderer
}

type Handler interface {
	Handle(ctx context.Context, messages ...interface{})
}
