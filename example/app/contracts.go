package app

import "errors"

type ProcessPaymentCommand struct {
	AccountID       uint64
	PaymentMethodID uint64
	OrderID         uint64
	Amount          uint64
	UserAgent       string
	Result          struct {
		PaymentID uint64
		Error     error
	}
}

var ErrPaymentMethodDeclined = errors.New("payment declined")
