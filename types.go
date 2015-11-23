package binding

import "net/http"

type (
	InputFactory     func() interface{}
	ControllerAction func(interface{}) Renderer

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
