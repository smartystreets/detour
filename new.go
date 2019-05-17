package detour

import (
	"fmt"
	"net/http"
	"reflect"
)

type (
	createModel func() interface{}
	monadicAction func(interface{}) Renderer
	niladicAction func() Renderer
)

func NewFromFactory(inputModelFactory createModel, controllerAction interface{}) http.Handler {
	expectedModelType := identifyInputModelArgumentType(controllerAction)
	if expectedModelType == nil {
		panic("Controller action must accept an input model.")
	}

	actualModelType := reflect.TypeOf(inputModelFactory())
	if actualModelType != expectedModelType {
		panic(fmt.Sprintf(
			"Controller requires input model of type: [%v] Factory function provided input model of type: [%v]",
			expectedModelType,
			actualModelType,
		))
	}

	return withFactory(controllerAction, inputModelFactory)
}

func New(controllerAction interface{}) http.Handler {
	modelType := identifyInputModelArgumentType(controllerAction)
	if modelType == nil {
		return simple(controllerAction.(func() Renderer))
	}

	return withFactory(controllerAction, func() interface{} {
		return reflect.New(modelType.Elem()).Interface()
	})
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
	if !isMethod(actionType) {
		panic("The action provided is not a function.")
	}

	if !returnsRenderer(actionType) {
		panic("The return type must implement the detour.Renderer interface.")
	}

	argumentCount := actionType.NumIn()
	if argumentCount == 0 {
		return nil
	}

	if argumentCount > 1 {
		panic("The callback provided must have no more than one argument.")
	}

	firstArgumentType := actionType.In(0)
	if !isSinglePointerArgument(argumentCount, firstArgumentType) {
		panic("The first argument to the controller callback must be a pointer type.")
	}

	return firstArgumentType
}

func isMethod(callback reflect.Type) bool {
	return callback.Kind() == reflect.Func
}

func returnsRenderer(actionType reflect.Type) bool {
	return actionType.NumOut() == 1 && actionType.Out(0).Implements(renderer)
}

var renderer = reflect.TypeOf((*Renderer)(nil)).Elem()

func isSinglePointerArgument(argumentCount int, firstArgumentType reflect.Type) bool {
	return argumentCount == 1 && isPointer(firstArgumentType)
}

func isPointer(argumentType reflect.Type) bool {
	return argumentType.Kind() == reflect.Ptr
}
