package ordersaga

import (
	"context"

	"github.com/Azure/go-shuttle/message"

	"saga-example/messages"
)

func (s *Saga) HandlePaymentVerified(ctx context.Context, paymentVerified *messages.PaymentVerified) message.Handler {
	// payment verified logic
	return paymentVerified.Complete()
}
