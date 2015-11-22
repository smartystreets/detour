package binding

import (
	"fmt"
	"net/http"
	"reflect"
)

type ModelBinder struct {
	domain     DomainAction
	controller ControllerAction
	input      InputFactory
}

func Typed(controllerAction interface{}) *ModelBinder {
	inputType := parseInputModelType(controllerAction).Elem()
	var factory InputFactory = func() interface{} { return reflect.New(inputType).Interface() }
	return TypedFactory(controllerAction, factory)
}

func TypedFactory(controllerAction interface{}, input InputFactory) *ModelBinder {
	callbackType := reflect.ValueOf(controllerAction)
	var callback ControllerAction = func(w http.ResponseWriter, r *http.Request, m interface{}) {
		callbackType.Call([]reflect.Value{reflect.ValueOf(w), reflect.ValueOf(r), reflect.ValueOf(m)})
	}
	return GenericFactory(callback, input)
}

func Generic(callback ControllerAction, message interface{}) *ModelBinder {
	inputType := reflect.TypeOf(message).Elem()
	var factory InputFactory = func() interface{} { return reflect.New(inputType).Interface() }
	return GenericFactory(callback, factory)
}

func GenericFactory(callback ControllerAction, input InputFactory) *ModelBinder {
	return &ModelBinder{controller: callback, input: input}
}

func Domain(callback DomainAction, input InputFactory) *ModelBinder {
	return &ModelBinder{
		domain: callback,
		input:  input,
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
	// FUTURE: if request has a Body (PUT/POST) and Content-Type: application/json
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
}

//////////////////////////////////////////////////////////////////////////////

const httpStatusUnprocessableEntity = 422
