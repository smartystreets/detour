package main

import (
	"net/http"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/detour/v3/detourtest"
	"github.com/smartystreets/detour/v3/example/app"
	"github.com/smartystreets/gunit"
)

func TestProcessPaymentDetourFixture(t *testing.T) {
	gunit.Run(new(ProcessPaymentDetourFixture), t)
}

type ProcessPaymentDetourFixture struct {
	*gunit.Fixture
	*detourtest.DetourFixture
}

func (this *ProcessPaymentDetourFixture) Setup() {
	this.DetourFixture = detourtest.Initialize()
	this.RequestHeaders.Set("Account-Id", "1")
	this.RequestHeaders.Set("User-Agent", "UserAgent")
	this.RequestBody["amount"] = 2
	this.RequestBody["order_id"] = 3
	this.RequestBody["payment_method_id"] = 4
}
func (this *ProcessPaymentDetourFixture) Teardown() {
	if this.Failed() {
		this.Print(this.Dump.String())
	}
}

func (this *ProcessPaymentDetourFixture) Test_NoAccountID_HTTP500() {
	this.RequestHeaders.Del("Account-Id")

	this.Do(NewProcessPaymentDetour)

	this.So(this.Handler.HandleCount, should.Equal, 0)
	this.So(this.ResponseStatus, should.Equal, http.StatusInternalServerError)
	this.So(this.ResponseBody, should.BeBlank)
}
func (this *ProcessPaymentDetourFixture) Test_MalformedJSON_HTTP400() {
	this.RequestBody["invalid-type"] = make(chan int)

	this.Do(NewProcessPaymentDetour)

	this.So(this.Handler.HandleCount, should.Equal, 0)
	this.So(this.ResponseStatus, should.Equal, http.StatusBadRequest)
	this.So(this.ResponseBody, should.BeBlank)
}
func (this *ProcessPaymentDetourFixture) Test_OrderIDRequired_HTTP422() {
	delete(this.RequestBody, "order_id")

	this.Do(NewProcessPaymentDetour)

	this.So(this.Handler.HandleCount, should.Equal, 0)
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
	delete(this.RequestBody, "amount")

	this.Do(NewProcessPaymentDetour)

	this.So(this.Handler.HandleCount, should.Equal, 0)
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
	delete(this.RequestBody, "payment_method_id")

	this.Do(NewProcessPaymentDetour)

	this.So(this.Handler.HandleCount, should.Equal, 0)
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
	delete(this.RequestBody, "amount")
	delete(this.RequestBody, "order_id")
	delete(this.RequestBody, "payment_method_id")

	this.Do(NewProcessPaymentDetour)

	this.So(this.Handler.HandleCount, should.Equal, 0)
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
	this.Handler.Prepare(&app.ProcessPaymentCommand{}, func(command interface{}) {
		command.(*app.ProcessPaymentCommand).Result.Error = app.ErrPaymentMethodDeclined
	})

	this.Do(NewProcessPaymentDetour)

	this.So(this.ResponseStatus, should.Equal, http.StatusPaymentRequired)
	this.So(this.ResponseHeaders.Get("Content-Type"), should.Equal, "application/json; charset=utf-8")
	this.So(this.ResponseBodyJSON(), should.Resemble, map[string]interface{}{
		"result":  "payment-context:payment-method-declined",
		"message": "The payment method was declined by the issuing bank.",
	})
}
func (this *ProcessPaymentDetourFixture) Test_HAPPY_HTTP200() {
	this.Handler.Prepare(&app.ProcessPaymentCommand{}, func(command interface{}) {
		command.(*app.ProcessPaymentCommand).Result.PaymentID = 42
	})

	this.Do(NewProcessPaymentDetour)

	this.So(this.ResponseStatus, should.Equal, http.StatusOK)
	this.So(this.ResponseHeaders.Get("Content-Type"), should.Equal, "application/json; charset=utf-8")
	this.So(this.ResponseBodyJSON(), should.Resemble, map[string]interface{}{
		"result": "payment-context:ok",
		"data": map[string]interface{}{
			"payment_id": 42.0,
		},
	})
	this.So(this.Handler.Context, should.Equal, this.RequestContext)
	this.So(this.Handler.Messages, should.Resemble, []interface{}{
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
