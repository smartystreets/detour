package detour

import (
	"net/http"
	"reflect"
)

type (
	createModel   func() interface{}
	monadicAction func(interface{}) Renderer
	niladicAction func() Renderer
)

func New(controllerAction interface{}) http.Handler {
	modelType := identifyInputModelArgumentType(controllerAction)
	if modelType == nil {
		return simple(controllerAction.(func() Renderer))
	}

	modelElement := modelType.Elem() // do not inline into factory callback method
	var factory createModel = func() interface{} { return reflect.New(modelElement).Interface() }
	return withFactory(controllerAction, factory)
}

func New2(controllerAction interface{}) http.Handler {
	modelType := identifyInputModelArgumentType(controllerAction)
	if modelType == nil {
		return simple(controllerAction.(func() Renderer))
	}

	var factory createModel = func() interface{} { return reflect.New(modelType.Elem()).Interface() }
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
