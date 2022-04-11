package shipping

import (
	"context"

	"saga-example/bus"
	"saga-example/messages"

	"github.com/Azure/go-shuttle/message"
)

type ShipOrderHandler struct {
	bus bus.Bus
}

func (p *ShipOrderHandler) HandleShipOrder(ctx context.Context, ship messages.ShipOrder) message.Handler {
	// only do work if order needs shipping (idempotent)
	// always publish event, republishing is ok.
	if err := p.bus.PublishEvent(ctx, &messages.OrderShipped{OrderID: ship.OrderID}); err != nil {
		return ship.Error(err)
	}
	return ship.Complete()
}
