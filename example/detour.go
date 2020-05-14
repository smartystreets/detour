package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/smartystreets/detour/v3"
	"github.com/smartystreets/detour/v3/example/app"
	"github.com/smartystreets/detour/v3/render"
)

type ProcessPaymentDetour struct {
	accountID       uint64
	PaymentMethodID uint64 `json:"payment_method_id"`
	OrderID         uint64 `json:"order_id"`
	Amount          uint64 `json:"amount"`

	command  *app.ProcessPaymentCommand
	renderer render.Renderer
}

func NewProcessPaymentDetour() detour.Detour {
	return &ProcessPaymentDetour{}
}

func (this *ProcessPaymentDetour) Bind(request *http.Request) []interface{} {
	this.bind(request)
	if this.renderer != nil {
		return nil
	}

	this.validate()
	if this.renderer != nil {
		return nil
	}

	return this.appendMessage(request)
}

func (this *ProcessPaymentDetour) bind(request *http.Request) {
	this.accountID, _ = strconv.ParseUint(request.Header.Get("Account-Id"), 10, 64)
	if this.accountID == 0 {
		this.renderer = render.StatusCodeResult{StatusCode: http.StatusInternalServerError}
		return
	}

	err := json.NewDecoder(request.Body).Decode(this)
	if err != nil {
		this.renderer = render.StatusCodeResult{StatusCode: http.StatusBadRequest}
	}
}

func (this *ProcessPaymentDetour) validate() {
	var validation []FieldError
	if this.Amount == 0 {
		validation = append(validation, invalidAmount)
	}
	if this.OrderID == 0 {
		validation = append(validation, invalidOrderID)
	}
	if len(validation) == 0 {
		return
	}

	result := LookupResult(errValidation)
	result.Data = validation
	this.renderer = render.JSONResult{StatusCode: result.StatusCode, Content: result}
}

func (this *ProcessPaymentDetour) appendMessage(request *http.Request) (messages []interface{}) {
	this.command = &app.ProcessPaymentCommand{
		AccountID:       this.accountID,
		Amount:          this.Amount,
		PaymentMethodID: this.PaymentMethodID,
		OrderID:         this.OrderID,
		UserAgent:       strings.TrimSpace(request.UserAgent()),
		UserAddress:     request.RemoteAddr,
	}
	return append(messages, this.command)
}

func (this *ProcessPaymentDetour) Render(response http.ResponseWriter, request *http.Request) {
	if this.renderer == nil {
		result := LookupResult(this.command.Result.Error)
		result.Data = ProcessedPaymentResult{PaymentID: this.command.Result.PaymentID}
		this.renderer = render.JSONResult{Content: result, Indent: "  "}
	}
	this.renderer.Render(response, request)
}

type ProcessedPaymentResult struct {
	PaymentID uint64 `json:"payment_id,omitempty"`
}
