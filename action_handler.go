package binding

import (
	"fmt"
	"net/http"
	"reflect"
)

type ActionHandler struct {
	controller ControllerAction
	input      CreateModel
}

func New(controllerAction interface{}) http.Handler {
	inputType := parseModelType(controllerAction).Elem()
	var factory CreateModel = func() interface{} { return reflect.New(inputType).Interface() }
	return withFactory(controllerAction, factory)
}

func withFactory(controllerAction interface{}, input CreateModel) *ActionHandler {
	callbackType := reflect.ValueOf(controllerAction)
	var callback ControllerAction = func(m interface{}) Renderer {
		results := callbackType.Call([]reflect.Value{reflect.ValueOf(m)})
		result := results[0]
		if result.IsNil() {
			return nil
		}
		return result.Elem().Interface().(Renderer)
	}
	return &ActionHandler{controller: callback, input: input}
}

func parseModelType(action interface{}) reflect.Type {
	actionType := reflect.TypeOf(action)
	if actionType.Kind() != reflect.Func {
		panic("The action provided is not a function.")
	} else if argumentCount := actionType.NumIn(); argumentCount != 1 {
		panic("The callback provided must have exactly one argument.")
	} else if actionType.In(0).Kind() != reflect.Ptr {
		panic("The first argument to the controller callback must be a pointer type.")
	} else if actionType.NumOut() != 1 || !actionType.Out(0).Implements(reflect.TypeOf((*Renderer)(nil)).Elem()) {
		panic("The return type must implement Renderer")
	} else {
		return actionType.In(0)
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
	response.WriteHeader(code)
	response.Header().Set("Content-Type", "application/json")
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
