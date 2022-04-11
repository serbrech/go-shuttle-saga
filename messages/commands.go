package messages

import "github.com/Azure/go-shuttle/message"

type CreateOrderCommand struct {
	*message.Message
	OrderID string
}

type CompleteOrderCommand struct {
	*message.Message
	OrderID string
}

type ShipOrder struct {
	*message.Message
	OrderID string
}

type ProcessPayment struct {
	*message.Message
	OrderID string
}
