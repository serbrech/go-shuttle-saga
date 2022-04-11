package ordersaga

import (
	"context"

	"github.com/Azure/go-shuttle/message"

	"saga-example/messages"
)

func (s *Saga) HandleOrderShipped(ctx context.Context, orderShipped *messages.OrderShipped) message.Handler {
	s.State.IsOrderShipped = true
	if s.IsOrderSagaCompleted() {
		if err := s.bus.SendCommand(ctx, messages.CompleteOrderCommand{OrderID: s.State.OrderID}); err != nil {
			return orderShipped.Error(err)
		}
	}
	return orderShipped.Complete()
}
