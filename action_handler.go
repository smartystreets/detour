package binding

import (
	"fmt"
	"net/http"
	"reflect"
)

type ActionHandler struct {
	controller MonadicAction
	input      CreateModel
}

func New(controllerAction interface{}) http.Handler {
	modelType := parseModelType(controllerAction)
	if modelType == nil {
		return simple(controllerAction.(func() Renderer))
	}

	modelElement := modelType.Elem() // do not inline into factory callback method
	var factory CreateModel = func() interface{} { return reflect.New(modelElement).Interface() }
	return withFactory(controllerAction, factory)
}

func withFactory(controllerAction interface{}, input CreateModel) http.Handler {
	callbackType := reflect.ValueOf(controllerAction)
	var callback MonadicAction = func(m interface{}) Renderer {
		results := callbackType.Call([]reflect.Value{reflect.ValueOf(m)})
		result := results[0]
		if result.IsNil() {
			return nil
		}
		return result.Elem().Interface().(Renderer)
	}
	return &ActionHandler{controller: callback, input: input}
}

func simple(controllerAction NiladicAction) http.Handler {
	return &ActionHandler{
		controller: func(interface{}) Renderer { return controllerAction() },
		input:      func() interface{} { return nil },
	}
}

func parseModelType(action interface{}) reflect.Type {
	actionType := reflect.TypeOf(action)
	if actionType.Kind() != reflect.Func {
		panic("The action provided is not a function.")
	} else if argumentCount := actionType.NumIn(); argumentCount > 1 {
		panic("The callback provided must have no more than one argument.")
	} else if argumentCount > 0 && actionType.In(0).Kind() != reflect.Ptr {
		panic("The first argument to the controller callback must be a pointer type.")
	} else if actionType.NumOut() != 1 || !actionType.Out(0).Implements(reflect.TypeOf((*Renderer)(nil)).Elem()) {
		panic("The return type must implement Renderer")
	} else if argumentCount > 0 {
		return actionType.In(0)
	} else {
		return nil
	}
}

//////////////////////////////////////////////////////////////////////////////

func (this *ActionHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
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
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(code)
	fmt.Fprint(response, err.Error())
}

func (this *ActionHandler) bind(request *http.Request, message interface{}) error {
	// FUTURE: if request has a Body (PUT/POST) and Content-Type: application/json
	if binder, ok := message.(Binder); !ok {
		return nil
	} else if err := request.ParseForm(); err != nil {
		return err
	} else {
		return binder.Bind(request)
	}
}

func (this *ActionHandler) validate(message interface{}) error {
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

func (this *ActionHandler) handle(response http.ResponseWriter, request *http.Request, message interface{}) {
	if result := this.controller(message); result != nil {
		result.Render(response, request)
	}
}

const httpStatusUnprocessableEntity = 422
