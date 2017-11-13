package detour

import (
	"errors"
	"net/http"
	"reflect"
)

type actionHandler struct {
	controller            monadicAction
	generateNewInputModel createModel
}

func New(controllerAction interface{}) http.Handler {
	modelType := identifyInputModelArgumentType(controllerAction)
	if modelType == nil {
		return simple(controllerAction.(func() Renderer))
	}

	modelElement := modelType.Elem() // do not inline into factory callback method
	var factory createModel = func() interface{} { return reflect.New(modelElement).Interface() }
	return withFactory(controllerAction, factory)
}

func withFactory(controllerAction interface{}, input createModel) http.Handler {
	callbackType := reflect.ValueOf(controllerAction)
	var callback monadicAction = func(m interface{}) Renderer {
		results := callbackType.Call([]reflect.Value{reflect.ValueOf(m)})
		result := results[0]
		if result.IsNil() {
			return nil
		}
		return result.Elem().Interface().(Renderer)
	}
	return &actionHandler{controller: callback, generateNewInputModel: input}
}

func simple(controllerAction niladicAction) http.Handler {
	return &actionHandler{
		controller:            func(interface{}) Renderer { return controllerAction() },
		generateNewInputModel: func() interface{} { return nil },
	}
}

func identifyInputModelArgumentType(action interface{}) reflect.Type {
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

// Install merely allows *actionHandler to implement a non-public/internal, company-specific interface.
func (this *actionHandler) Install(http.Handler) {}

func (this *actionHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	model := this.generateNewInputModel()
	status, err := prepareInputModel(model, request)
	result := this.determineResult(model, status, err)
	result.Render(response, request)
}

func prepareInputModel(model interface{}, request *http.Request) (int, error) {
	if err := bind(request, model); err != nil {
		return http.StatusBadRequest, err
	}

	sanitize(model)

	if err := validate(model); err != nil {
		return http.StatusUnprocessableEntity, err
	}

	if err := serverError(model); err != nil {
		return http.StatusInternalServerError, err
	}

	return 0, nil
}

func bind(request *http.Request, message interface{}) error {
	// FUTURE: if request has a Body (PUT/POST) and Content-Type: application/json
	if binder, isBinder := message.(Binder); !isBinder {
		return nil
	} else if err := request.ParseForm(); err != nil {
		return err
	} else if err = binder.Bind(request); err == nil {
		return nil
	} else if errs, isErrors := err.(Errors); isErrors && len(errs) == 0 {
		return nil
	} else {
		return err
	}
}

func sanitize(message interface{}) {
	if sanitizer, isSanitizer := message.(Sanitizer); isSanitizer {
		sanitizer.Sanitize()
	}
}
func validate(message interface{}) error {
	if validator, isValidator := message.(Validator); !isValidator {
		return nil
	} else if err := validator.Validate(); err == nil {
		return nil
	} else if errs, isErrors := err.(Errors); isErrors && len(errs) == 0 {
		return nil
	} else {
		return err
	}
}

func serverError(message interface{}) error {
	if server, isServerError := message.(ServerError); isServerError && server.Error() {
		return internalServerError
	}
	return nil
}

var internalServerError = errors.New(http.StatusText(http.StatusInternalServerError))

func (this *actionHandler) determineResult(model interface{}, status int, err error) Renderer {
	if err != nil {
		return inputModelErrorResult(status, err)
	} else {
		return this.controllerActionResult(model)
	}
}

func inputModelErrorResult(code int, err error) Renderer {
	if _, isErrors := err.(Errors); isErrors {
		return &JSONResult{StatusCode: code, Content: err}
	} else if _, isDiagnosticErr := err.(*DiagnosticError); isDiagnosticErr {
		return &DiagnosticResult{StatusCode: code, Message: err.Error()}
	} else {
		return &StatusCodeResult{StatusCode: code, Message: err.Error()}
	}
}

func (this *actionHandler) controllerActionResult(model interface{}) Renderer {
	if result := this.controller(model); result != nil {
		return result
	} else {
		return nopResult{}
	}
}

type nopResult struct{}

func (nopResult) Render(http.ResponseWriter, *http.Request) {}
