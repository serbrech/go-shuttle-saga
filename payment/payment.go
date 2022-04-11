package payment

import (
	"context"

	"saga-example/bus"
	"saga-example/messages"

	"github.com/Azure/go-shuttle/message"
)

type ProcessPaymentHandler struct {
	bus bus.Bus
}

func (p *ProcessPaymentHandler) HandleCreateOrder(ctx context.Context, payment messages.ProcessPayment) message.Handler {
	// handle payment
	if err := p.bus.PublishEvent(ctx, &messages.PaymentVerified{OrderID: payment.OrderID}); err != nil {
		return payment.Error(err)
	}
	return payment.Complete()
}
