package detour

import (
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
	message := this.input()
	if !this.bind(request, message, response) {
		return
	}
	this.sanitize(message)

	if !this.validate(message, response, request) {
		return
	}
	if !this.error(message, response) {
		return
	}
	this.handle(message, response, request)
}

// Install merely allows *ActionHandler to implement a non-public/internal, company-specific interface.
func (this *ActionHandler) Install(http.Handler) {}

func (this *ActionHandler) bind(request *http.Request, message interface{}, response http.ResponseWriter) bool {
	if err := bind(request, message); err != nil {
		writeErrorResponse(response, request, err, http.StatusBadRequest)
		return false
	}
	return true
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

func (this *ActionHandler) sanitize(message interface{}) {
	if sanitizer, ok := message.(Sanitizer); ok {
		sanitizer.Sanitize()
	}
}
func (this *ActionHandler) validate(message interface{}, response http.ResponseWriter, request *http.Request) bool {
	if err := validate(message); err != nil {
		writeErrorResponse(response, request, err, http.StatusUnprocessableEntity)
		return false
	}
	return true
}
func validate(message interface{}) error {
	if validator, ok := message.(Validator); !ok {
		return nil
	} else if err := validator.Validate(); err == nil {
		return nil
	} else if errors, ok := err.(*Errors); ok && len(errors.errors) == 0 {
		return nil
	} else {
		return err
	}
}

func (this *ActionHandler) error(message interface{}, response http.ResponseWriter) bool {
	if server, ok := message.(ServerError); ok && server.Error() {
		http.Error(response, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return false
	}
	return true
}

func (this *ActionHandler) handle(message interface{}, response http.ResponseWriter, request *http.Request) {
	if result := this.controller(message); result != nil {
		result.Render(response, request)
	}
}

func writeErrorResponse(response http.ResponseWriter, request *http.Request, err error, code int) {
	var result Renderer

	if _, ok := err.(*Errors); ok {
		result = &JSONResult{StatusCode: code, Content: err}
	} else if _, ok := err.(*DiagnosticError); ok {
		result = &DiagnosticResult{StatusCode: code, Message: err.Error()}
	} else {
		result = &StatusCodeResult{StatusCode: code, Message: err.Error()}
	}

	result.Render(response, request)
}
