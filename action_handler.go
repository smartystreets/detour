package detour

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

func (this *ActionHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	NewAction(this.input(), this.controller, request, response).Execute()
}

///////////////////////////////////////////////////////////////////////////////

type Action struct {
	request    *http.Request
	response   http.ResponseWriter
	controller MonadicAction
	message    interface{}
	finished   bool
}

func NewAction(message interface{}, controller MonadicAction, request *http.Request, response http.ResponseWriter) *Action {
	return &Action{
		message:    message,
		controller: controller,
		request:    request,
		response:   response,
	}
}

func (this *Action) Execute() {
	this.step(this.bind)
	this.step(this.sanitize)
	this.step(this.validate)
	this.step(this.error)
	this.step(this.handle)
}
func (this *Action) step(action func()) {
	if !this.finished {
		action()
	}
}
func (this *Action) bind() {
	if err := bind(this.request, this.message); err != nil {
		writeJSONError(this.response, err, http.StatusBadRequest)
		this.finished = true
	}
}
func (this *Action) sanitize() {
	if sanitizer, ok := this.message.(Sanitizer); ok {
		sanitizer.Sanitize()
	}
}
func (this *Action) validate() {
	if err := validate(this.message); err != nil {
		writeJSONError(this.response, err, httpStatusUnprocessableEntity)
		this.finished = true
	}
}
func (this *Action) error() {
	if server, ok := this.message.(ServerError); ok && server.Error() {
		writeInternalServerError(this.response)
		this.finished = true
	}
}
func (this *Action) handle() {
	if result := this.controller(this.message); result != nil {
		result.Render(this.response, this.request)
	}
}

func bind(request *http.Request, message interface{}) error {
	// FUTURE: if request has a Body (PUT/POST) and Content-Type: application/json
	if binder, ok := message.(Binder); !ok {
		return nil
	} else if err := request.ParseForm(); err != nil {
		return err
	} else {
		return binder.Bind(request)
	}
}
func validate(message interface{}) error {
	if validator, ok := message.(Validator); !ok {
		return nil
	} else if err := validator.Validate(); err == nil {
		return nil
	} else if errors, ok := err.(Errors); ok && len(errors) == 0 {
		return nil
	} else {
		return err
	}
}
func writeJSONError(response http.ResponseWriter, err error, code int) {
	response.Header().Set(contentTypeHeader, jsonContentType)
	response.WriteHeader(code)
	fmt.Fprint(response, err.Error())
}
func writeInternalServerError(response http.ResponseWriter) {
	http.Error(response, internalServerErrorText, http.StatusInternalServerError)
}

const httpStatusUnprocessableEntity = 422

var internalServerErrorText = http.StatusText(http.StatusInternalServerError)
