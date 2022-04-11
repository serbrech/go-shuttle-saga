package messages

import "github.com/Azure/go-shuttle/message"

type PaymentVerified struct {
	*message.Message
	OrderID string
}

type OrderShipped struct {
	*message.Message
	OrderID string
}
