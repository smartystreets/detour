package binding

import (
	"fmt"
	"net/http"
	"reflect"
)

type ModelBinder struct {
	input      InputModelFactory
	domain     DomainAction
	controller ControllerAction
}

func DefaultBinder(controllerAction interface{}) *ModelBinder {
	inputType := parseInputModelType(controllerAction).Elem()
	callback := reflect.ValueOf(controllerAction)
	return ControllerBinder(
		func() interface{} { return reflect.New(inputType).Interface() },
		func(w http.ResponseWriter, r *http.Request, m interface{}) {
			callback.Call([]reflect.Value{reflect.ValueOf(w), reflect.ValueOf(r), reflect.ValueOf(m)})
		},
	)
}

func ControllerBinder(input InputModelFactory, callback ControllerAction) *ModelBinder {
	return &ModelBinder{
		input:      input,
		controller: callback,
	}
}

func DomainBinder(input InputModelFactory, callback DomainAction) *ModelBinder {
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
	} else if this.controller != nil {
		this.controller(response, request, message)
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
	// FUTURE: if request has a Body (PUT/POST)
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
	} else if err := validator.Validate(); err == nil {
		return nil
	} else if errors, ok := err.(ValidationErrors); ok && len(errors) == 0 {
		return nil
	} else {
		return err
	}
}

func (this *ModelBinder) handle(response http.ResponseWriter, request *http.Request, message interface{}) {
	if translator, ok := message.(Translator); ok {
		message = translator.Translate()
	}

	if result := this.domain(message); result != nil {
		result.ServeHTTP(response, request)
	}
}

//////////////////////////////////////////////////////////////////////////////

func parseInputModelType(function interface{}) reflect.Type {
	typed := reflect.TypeOf(function)
	if typed.Kind() != reflect.Func {
		panic("The controller callback provided is not a function.")
	} else if argumentCount := typed.NumIn(); argumentCount != 3 {
		panic("The controller callback provided must have exactly three arguments.")
	} else if !typed.In(0).Implements(reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()) {
		panic("The first argument to the controller callback must of type http.ResponseWriter.")
	} else if typed.In(1) != reflect.TypeOf(&http.Request{}) {
		panic("The second argument to the controller callback must of type http.ResponseWriter.")
	} else if typed.In(2).Kind() != reflect.Ptr {
		panic("The third argument to the controller callback must be a pointer type.")
	} else {
		return typed.In(2)
	}

	return reflect.TypeOf(0)
}

//////////////////////////////////////////////////////////////////////////////

const httpStatusUnprocessableEntity = 422
