package binding

import "net/http"

type (
	InputFactory     func() interface{}
	DomainAction     func(interface{}) http.Handler
	ControllerAction func(http.ResponseWriter, *http.Request, interface{})

	Binder interface {
		Bind(*http.Request) error
	}

	Validator interface {
		Validate() error
	}

	Translator interface {
		Translate() interface{}
	}
)
