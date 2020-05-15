package main

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/detour/v3"
	"github.com/smartystreets/detour/v3/example/app"
	"github.com/smartystreets/detour/v3/httptest"
	"github.com/smartystreets/gunit"
)

func TestProcessPaymentDetourFixture(t *testing.T) {
	gunit.Run(new(ProcessPaymentDetourFixture), t)
}

type ProcessPaymentDetourFixture struct {
	*gunit.Fixture
	*httptest.HTTPFixture
	handler *FakeHandler

	jsonRequestBody map[string]interface{}
}

func (this *ProcessPaymentDetourFixture) Setup() {
	this.handler = NewFakeHandler()
	this.HTTPFixture = httptest.NewHTTPFixture()
	this.RequestHeaders.Set("Account-Id", "1")
	this.RequestHeaders.Set("User-Agent", "UserAgent")
	this.jsonRequestBody = map[string]interface{}{
		"amount":            2,
		"order_id":          3,
		"payment_method_id": 4,
	}
}
func (this *ProcessPaymentDetourFixture) Detour() {
	this.SetJSONBody(this.jsonRequestBody)
	this.Serve(detour.New(NewProcessPaymentDetour, this.handler))
}
func (this *ProcessPaymentDetourFixture) Teardown() {
	if this.Failed() {
		this.Print(this.Dump.String())
	}
}

func (this *ProcessPaymentDetourFixture) Test_NoAccountID_HTTP500() {
	this.RequestHeaders.Del("Account-Id")

	this.Detour()

	this.So(this.handler.HandleCount, should.Equal, 0)
	this.So(this.ResponseStatus, should.Equal, http.StatusInternalServerError)
	this.So(this.ResponseBody, should.BeBlank)
}
func (this *ProcessPaymentDetourFixture) Test_MalformedJSON_HTTP400() {
	this.jsonRequestBody["invalid-type"] = make(chan int)

	this.Detour()

	this.So(this.handler.HandleCount, should.Equal, 0)
	this.So(this.ResponseStatus, should.Equal, http.StatusBadRequest)
	this.So(this.ResponseBody, should.BeBlank)
}
func (this *ProcessPaymentDetourFixture) Test_OrderIDRequired_HTTP422() {
	delete(this.jsonRequestBody, "order_id")

	this.Detour()

	this.So(this.handler.HandleCount, should.Equal, 0)
	this.So(this.ResponseStatus, should.Equal, http.StatusUnprocessableEntity)
	this.So(this.ResponseHeaders.Get("Content-Type"), should.Equal, "application/json; charset=utf-8")
	this.So(this.ResponseBodyJSON(), should.Resemble, map[string]interface{}{
		"result":  "payment-context:request-validation-error",
		"message": "The request was invalid. See included data for details.",
		"data": []interface{}{
			map[string]interface{}{
				"field_name": "order_id",
				"message":    "a valid order id is required",
			},
		},
	})
}
func (this *ProcessPaymentDetourFixture) Test_AmountRequired_HTTP422() {
	delete(this.jsonRequestBody, "amount")

	this.Detour()

	this.So(this.handler.HandleCount, should.Equal, 0)
	this.So(this.ResponseStatus, should.Equal, http.StatusUnprocessableEntity)
	this.So(this.ResponseHeaders.Get("Content-Type"), should.Equal, "application/json; charset=utf-8")
	this.So(this.ResponseBodyJSON(), should.Resemble, map[string]interface{}{
		"result":  "payment-context:request-validation-error",
		"message": "The request was invalid. See included data for details.",
		"data": []interface{}{
			map[string]interface{}{
				"field_name": "amount",
				"message":    "a positive amount is required",
			},
		},
	})
}
func (this *ProcessPaymentDetourFixture) Test_PaymentMethodIDRequired_HTTP422() {
	delete(this.jsonRequestBody, "payment_method_id")

	this.Detour()

	this.So(this.handler.HandleCount, should.Equal, 0)
	this.So(this.ResponseStatus, should.Equal, http.StatusUnprocessableEntity)
	this.So(this.ResponseHeaders.Get("Content-Type"), should.Equal, "application/json; charset=utf-8")
	this.So(this.ResponseBodyJSON(), should.Resemble, map[string]interface{}{
		"result":  "payment-context:request-validation-error",
		"message": "The request was invalid. See included data for details.",
		"data": []interface{}{
			map[string]interface{}{
				"field_name": "payment_method_id",
				"message":    "a valid payment method id is required",
			},
		},
	})
}
func (this *ProcessPaymentDetourFixture) Test_MultipleRequiredFields_HTTP422() {
	delete(this.jsonRequestBody, "amount")
	delete(this.jsonRequestBody, "order_id")
	delete(this.jsonRequestBody, "payment_method_id")

	this.Detour()

	this.So(this.handler.HandleCount, should.Equal, 0)
	this.So(this.ResponseStatus, should.Equal, http.StatusUnprocessableEntity)
	this.So(this.ResponseHeaders.Get("Content-Type"), should.Equal, "application/json; charset=utf-8")
	this.So(this.ResponseBodyJSON(), should.Resemble, map[string]interface{}{
		"result":  "payment-context:request-validation-error",
		"message": "The request was invalid. See included data for details.",
		"data": []interface{}{
			map[string]interface{}{
				"field_name": "amount",
				"message":    "a positive amount is required",
			},
			map[string]interface{}{
				"field_name": "order_id",
				"message":    "a valid order id is required",
			},
			map[string]interface{}{
				"field_name": "payment_method_id",
				"message":    "a valid payment method id is required",
			},
		},
	})
}
func (this *ProcessPaymentDetourFixture) Test_DeclineErrorFromApplication_HTTP402() {
	this.handler.Prepare(&app.ProcessPaymentCommand{}, func(command interface{}) {
		command.(*app.ProcessPaymentCommand).Result.Error = app.ErrPaymentMethodDeclined
	})

	this.Detour()

	this.So(this.ResponseStatus, should.Equal, http.StatusPaymentRequired)
	this.So(this.ResponseHeaders.Get("Content-Type"), should.Equal, "application/json; charset=utf-8")
	this.So(this.ResponseBodyJSON(), should.Resemble, map[string]interface{}{
		"result":  "payment-context:payment-method-declined",
		"message": "The payment method was declined by the issuing bank.",
	})
}
func (this *ProcessPaymentDetourFixture) Test_HAPPY_HTTP200() {
	this.handler.Prepare(&app.ProcessPaymentCommand{}, func(command interface{}) {
		command.(*app.ProcessPaymentCommand).Result.PaymentID = 42
	})

	this.Detour()

	this.So(this.ResponseStatus, should.Equal, http.StatusOK)
	this.So(this.ResponseHeaders.Get("Content-Type"), should.Equal, "application/json; charset=utf-8")
	this.So(this.ResponseBodyJSON(), should.Resemble, map[string]interface{}{
		"result": "payment-context:ok",
		"data": map[string]interface{}{
			"payment_id": 42.0,
		},
	})
	this.So(this.handler.Context, should.Equal, this.RequestContext)
	this.So(this.handler.Messages, should.Resemble, []interface{}{
		&app.ProcessPaymentCommand{
			AccountID:       1,
			Amount:          2,
			OrderID:         3,
			PaymentMethodID: 4,
			UserAgent:       "UserAgent",
			Result: struct {
				PaymentID uint64
				Error     error
			}{
				PaymentID: 42,
				Error:     nil,
			},
		},
	})
}

////////////////////////////////////////////////////////////////

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
