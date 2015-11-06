package binding

import (
	"fmt"
	"net/http"
)

type ModelBinder struct {
	input   InputModelFactory
	domain  DomainHandler
	handler Handler
}

func NewModelBinderHandler(input InputModelFactory, handler Handler) *ModelBinder {
	return &ModelBinder{
		input:   input,
		handler: handler,
	}
}

func NewDomainModelBinder(input InputModelFactory, domain DomainHandler) *ModelBinder {
	return &ModelBinder{
		input:  input,
		domain: domain,
	}
}

func (this *ModelBinder) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	message := this.input()

	if err := this.bind(request, message); err != nil {
		writeJSONError(response, err, http.StatusBadRequest)
	} else if err := this.validate(message); err != nil {
		writeJSONError(response, err, httpStatusUnprocessableEntity)
	} else {
		this.handle(response, request, message)
	}
}

func writeJSONError(response http.ResponseWriter, err error, code int) {
	response.WriteHeader(code)
	response.Header().Set("Content-Type", "application/json")
	fmt.Fprint(response, err.Error())
}

func (this *ModelBinder) bind(request *http.Request, message interface{}) error {
	if binder, ok := message.(Binder); !ok {
		return nil
	} else {
		return binder.Bind(request)
	}
}

func (this *ModelBinder) validate(message interface{}) error {
	if validator, ok := message.(Validator); !ok {
		return nil
	} else {
		return validator.Validate()
	}
}

func (this *ModelBinder) handle(response http.ResponseWriter, request *http.Request, message interface{}) {
	if this.handler != nil {
		this.handler.Handle(response, request, message)
		return
	}

	if translator, ok := message.(Translator); ok {
		message = translator.Translate()
	}

	if result := this.domain(message); result != nil {
		result.ServeHTTP(response, request)
	}
}

//////////////////////////////////////////////////////////////////////////////

const httpStatusUnprocessableEntity = 422
