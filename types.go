package binding

import "net/http"

type (
	InputModelFactory func() interface{}

	DomainAction     func(interface{}) http.Handler
	ControllerAction func(response http.ResponseWriter, request *http.Request, message interface{})

	Binder interface {
		Bind(request *http.Request) error
	}

	Validator interface {
		Validate() error
	}

	Translator interface {
		Translate() interface{}
	}
)
