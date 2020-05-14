package detourtest

import (
	"context"
	"reflect"
)

type FakeHandler struct {
	HandleCount int
	Context     context.Context
	Messages    []interface{}
	Handlers    map[reflect.Type]func(interface{})
}

func NewFakeHandler() *FakeHandler {
	return &FakeHandler{Handlers: make(map[reflect.Type]func(interface{}))}
}

func (this *FakeHandler) Prepare(message interface{}, callback func(interface{})) {
	this.Handlers[reflect.TypeOf(message)] = callback
}

func (this *FakeHandler) Handle(ctx context.Context, messages ...interface{}) {
	this.HandleCount++
	this.Context = ctx
	this.Messages = messages

	for _, message := range messages {
		callback := this.Handlers[reflect.TypeOf(message)]
		if callback != nil {
			callback(message)
		}
	}
}
