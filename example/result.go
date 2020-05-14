package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/smartystreets/detour/v3/example/app"
)

type BasicResult struct {
	StatusCode int         `json:"-"`
	Result     string      `json:"result,omitempty"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data"`
}

func LookupResult(err error) BasicResult {
	result, found := results[err]
	if !found {
		result.StatusCode = http.StatusUnprocessableEntity
		result.Result = "payment-context:unrecognized-error"
		result.Message = fmt.Sprintf("Unrecognized error: [%s]", err)
	}
	return result
}

var results = map[error]BasicResult{
	nil: {
		StatusCode: http.StatusOK,
		Result:     "payment-context:ok",
		Message:    "",
	},
	errValidation: {
		StatusCode: http.StatusUnprocessableEntity,
		Result:     "payment-context:request-validation-error",
		Message:    "The request was invalid. See included data for details.",
	},
	app.ErrPaymentMethodDeclined: {
		StatusCode: http.StatusPaymentRequired,
		Result:     "payment-context:payment-method-declined",
		Message:    "The payment method was declined by the issuing bank.",
	},
}

var errValidation = errors.New("validation")

type FieldError struct {
	FieldName string `json:"field_name"`
	Message   string `json:"message"`
}

var (
	invalidAmount  = FieldError{FieldName: "amount", Message: "a positive amount is required"}
	invalidOrderID = FieldError{FieldName: "order_id", Message: "a valid order id is required"}
)
