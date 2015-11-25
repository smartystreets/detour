package detour

import "net/http"

type (
	CreateModel   func() interface{}
	MonadicAction func(interface{}) Renderer
	NiladicAction func() Renderer

	Binder interface {
		Bind(*http.Request) error
	}

	Validator interface {
		Validate() error
	}

	Renderer interface {
		Render(http.ResponseWriter, *http.Request)
	}
)
