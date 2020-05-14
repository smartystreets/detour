package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/detour/v3"
	"github.com/smartystreets/detour/v3/example/app"
	"github.com/smartystreets/gunit"
)

func TestProcessPaymentDetourFixture(t *testing.T) {
	gunit.Run(new(ProcessPaymentDetourFixture), t)
}

type ProcessPaymentDetourFixture struct {
	*gunit.Fixture

	Handler *FakeHandler

	RequestURL      url.URL
	RequestURLQuery url.Values
	RequestBody     map[string]interface{}
	RequestHeaders  http.Header
	RequestContext  context.Context

	ResponseStatus  int
	ResponseHeaders http.Header
	ResponseBody    string
}

func (this *ProcessPaymentDetourFixture) Setup() {
	this.Handler = NewFakeHandler()
	this.RequestHeaders = http.Header{
		"Account-Id": []string{"1"},
		"User-Agent": []string{"UserAgent"},
	}
	this.RequestBody = map[string]interface{}{
		"amount":            2,
		"order_id":          3,
		"payment_method_id": 4,
	}
}

func (this *ProcessPaymentDetourFixture) detour() {
	request := this.buildRequest()
	handler := detour.New(NewProcessPaymentDetour, this.Handler)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	this.collectResponse(recorder)
}
func (this *ProcessPaymentDetourFixture) buildRequest() *http.Request {
	body, _ := json.Marshal(this.RequestBody)
	request := httptest.NewRequest("GET", "/", bytes.NewReader(body))
	request.Header = this.RequestHeaders
	requestDump, _ := httputil.DumpRequest(request, true)
	this.Printf("REQUEST DUMP:\n%s\n\n", formatDump(">", string(requestDump)))
	this.RequestContext = request.Context()
	return request
}
func (this *ProcessPaymentDetourFixture) collectResponse(recorder *httptest.ResponseRecorder) {
	response := recorder.Result()
	responseDump, _ := httputil.DumpResponse(response, true)
	this.Printf("RESPONSE DUMP:\n%s\n\n", formatDump("<", string(responseDump)))
	body, _ := ioutil.ReadAll(response.Body)
	this.ResponseBody = string(body)
	this.ResponseStatus = response.StatusCode
	this.ResponseHeaders = response.Header
}
func formatDump(prefix, dump string) string {
	prefix = "\n" + prefix + " "
	lines := strings.Split(strings.TrimSpace(dump), "\n")
	return prefix + strings.Join(lines, prefix)
}

func (this *ProcessPaymentDetourFixture) ResponseBodyJSON() (actual map[string]interface{}) {
	this.So(this.ResponseHeaders.Get("Content-Type"), should.Equal, "application/json; charset=utf-8")
	err := json.Unmarshal([]byte(this.ResponseBody), &actual)
	this.So(err, should.BeNil)
	return actual
}

func (this *ProcessPaymentDetourFixture) Test_NoAccountID_HTTP500() {
	this.RequestHeaders.Del("Account-Id")

	this.detour()

	this.So(this.Handler.HandleCount, should.Equal, 0)
	this.So(this.ResponseStatus, should.Equal, http.StatusInternalServerError)
	this.So(this.ResponseBody, should.BeBlank)
}
func (this *ProcessPaymentDetourFixture) Test_MalformedJSON_HTTP400() {
	this.RequestBody["invalid-type"] = make(chan int)

	this.detour()

	this.So(this.Handler.HandleCount, should.Equal, 0)
	this.So(this.ResponseStatus, should.Equal, http.StatusBadRequest)
	this.So(this.ResponseBody, should.BeBlank)
}
func (this *ProcessPaymentDetourFixture) Test_OrderIDRequired_HTTP422() {
	delete(this.RequestBody, "order_id")

	this.detour()

	this.So(this.Handler.HandleCount, should.Equal, 0)
	this.So(this.ResponseStatus, should.Equal, http.StatusUnprocessableEntity)
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

	this.detour()

	this.So(this.Handler.HandleCount, should.Equal, 0)
	this.So(this.ResponseStatus, should.Equal, http.StatusUnprocessableEntity)
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

	this.detour()

	this.So(this.Handler.HandleCount, should.Equal, 0)
	this.So(this.ResponseStatus, should.Equal, http.StatusUnprocessableEntity)
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

	this.detour()

	this.So(this.Handler.HandleCount, should.Equal, 0)
	this.So(this.ResponseStatus, should.Equal, http.StatusUnprocessableEntity)
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

	this.detour()

	this.So(this.ResponseStatus, should.Equal, http.StatusPaymentRequired)
	this.So(this.ResponseBodyJSON(), should.Resemble, map[string]interface{}{
		"result":  "payment-context:payment-method-declined",
		"message": "The payment method was declined by the issuing bank.",
	})
}
func (this *ProcessPaymentDetourFixture) Test_HAPPY_HTTP200() {
	this.Handler.Prepare(&app.ProcessPaymentCommand{}, func(command interface{}) {
		command.(*app.ProcessPaymentCommand).Result.PaymentID = 42
	})

	this.detour()

	this.So(this.ResponseStatus, should.Equal, http.StatusOK)
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

/////////////////////////////////////////////////////////////////////

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
