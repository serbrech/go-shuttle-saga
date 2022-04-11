package ordersaga

import (
	"context"

	"github.com/Azure/go-shuttle/message"

	"saga-example/messages"
)

func (s *Saga) HandleCompleteOrder(ctx context.Context, completeOrder *messages.CompleteOrderCommand) message.Handler {
	// do some handling, maybe signal the completion of an operation.
	s.State.IsOrderCompleted = true
	if err := s.SaveState(); err != nil {
		return completeOrder.Error(err)
	}
	return completeOrder.Complete()
}
