package binding

import "net/http"

type (
	InputModelFactory func() interface{}

	DomainHandler func(interface{}) http.Handler

	Binder interface {
		Bind(request *http.Request) error
	}

	Validator interface {
		Validate() error
	}

	Translator interface {
		Translate() interface{}
	}

	Handler interface {
		Handle(response http.ResponseWriter, request *http.Request, message interface{})
	}
)
