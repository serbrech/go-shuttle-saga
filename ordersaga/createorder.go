package ordersaga

import (
	"context"
	"time"

	"github.com/Azure/go-shuttle/message"

	"saga-example/messages"
)

func (s *Saga) HandleCreateOrder(ctx context.Context, createOrder *messages.CreateOrderCommand) message.Handler {
	startTime := time.Now()
	if err := s.bus.SendCommand(ctx, &messages.ShipOrder{OrderID: createOrder.OrderID}); err != nil {
		return createOrder.Error(err)
	}
	if err := s.bus.SendCommand(ctx, &messages.ProcessPayment{OrderID: createOrder.OrderID}); err != nil {
		return createOrder.Error(err)
	}
	s.State.StartedAt = &startTime
	if err := s.SaveState(); err != nil {
		return createOrder.Error(err)
	}
	return createOrder.Complete()
}
