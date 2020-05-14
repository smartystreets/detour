package app

import "context"

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (this *Handler) Handle(_ context.Context, messages ...interface{}) {
	for _, message := range messages {
		switch message := message.(type) {
		case *ProcessPaymentCommand:
			message.Result.PaymentID = message.OrderID + message.Amount
		}
	}
}
