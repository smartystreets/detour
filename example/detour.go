package main

import (
	"context"
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
	userAgent       string
	userAddress     string
}

func NewProcessPaymentDetour() detour.Detour {
	return &ProcessPaymentDetour{}
}

func (this *ProcessPaymentDetour) Bind(request *http.Request) render.Renderer {
	this.accountID, _ = strconv.ParseUint(request.Header.Get("Account-Id"), 10, 64)
	if this.accountID == 0 {
		return render.StatusCodeResult{StatusCode: http.StatusInternalServerError}
	}

	err := json.NewDecoder(request.Body).Decode(this)
	if err != nil {
		return render.StatusCodeResult{StatusCode: http.StatusBadRequest}
	}

	var validation []FieldError
	if this.Amount == 0 {
		validation = append(validation, invalidAmount)
	}
	if this.OrderID == 0 {
		validation = append(validation, invalidOrderID)
	}
	if len(validation) > 0 {
		result := LookupResult(errValidation)
		result.Data = validation
		return render.JSONResult{StatusCode: result.StatusCode, Content: result}
	}

	this.userAgent = strings.TrimSpace(request.UserAgent())
	this.userAddress = request.RemoteAddr
	return nil
}

func (this *ProcessPaymentDetour) Handle(ctx context.Context, handler detour.Handler) render.Renderer {
	command := &app.ProcessPaymentCommand{
		AccountID:       this.accountID,
		Amount:          this.Amount,
		PaymentMethodID: this.PaymentMethodID,
		OrderID:         this.OrderID,
		UserAgent:       this.userAgent,
		UserAddress:     this.userAddress,
	}

	handler.Handle(ctx, command)

	result := LookupResult(command.Result.Error)
	result.Data = ProcessedPaymentResult{PaymentID: command.Result.PaymentID}
	return render.JSONResult{Content: result, Indent: "  "}
}

type ProcessedPaymentResult struct {
	PaymentID uint64 `json:"payment_id,omitempty"`
}
