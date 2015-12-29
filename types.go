package detour

import "net/http"

type (
	CreateModel   func() interface{}
	MonadicAction func(interface{}) Renderer
	NiladicAction func() Renderer

	Binder interface {
		Bind(*http.Request) error
	}

	Sanitizer interface {
		Sanitize()
	}

	Validator interface {
		Validate() error
	}

	ServerError interface {
		Error() bool
	}

	Renderer interface {
		Render(http.ResponseWriter, *http.Request)
	}
)
