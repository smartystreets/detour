package binding

import (
	"fmt"
	"net/http"
)

type ModelBinder struct {
	input      InputModelFactory
	domain     DomainAction
	controller ControllerAction
}

func NewBinder(input InputModelFactory, callback ControllerAction) *ModelBinder {
	return &ModelBinder{
		input:      input,
		controller: callback,
	}
}

func NewDomainBinder(input InputModelFactory, callback DomainAction) *ModelBinder {
	return &ModelBinder{
		input:  input,
		domain: callback,
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
	} else if err := request.ParseForm(); err != nil {
		return err
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
	if this.controller != nil {
		this.controller(response, request, message)
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
