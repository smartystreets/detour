package detour

import "net/http"

type (
	createModel   func() interface{}
	monadicAction func(interface{}) Renderer
	niladicAction func() Renderer
)

type (
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
