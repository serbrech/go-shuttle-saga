package ordersaga

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/go-shuttle/message"
	"github.com/Azure/go-shuttle/queue"
	ql "github.com/Azure/go-shuttle/queue/listener"
	"github.com/Azure/go-shuttle/topic"
	tl "github.com/Azure/go-shuttle/topic/listener"
	"github.com/prometheus/common/log"

	"saga-example/messages"
)

var (
	sagaName = "ordersaga"
)

func StartListeners(ctx context.Context) error {
	paymentEventListener, err := topic.NewListener(tl.WithSubscriptionName(sagaName))
	if err != nil {
		return err
	}
	listeners = append(listeners, paymentEventListener)

	shippingEventListener, err := topic.NewListener(tl.WithSubscriptionName(sagaName))
	if err != nil {
		return err
	}
	listeners = append(listeners, shippingEventListener)

	orderQueueListener, err := queue.NewListener(ql.WithMaxDeliveryCount(10))
	if err != nil {
		return err
	}
	listeners = append(listeners, orderQueueListener)

	go func() {
		drop := 0
		for err := paymentEventListener.Listen(ctx, message.HandleFunc(hydrateSagaAndDispatch), "payments"); err != nil; drop++ {
			log.Warnf("listener exited: %s", err)
		}
	}()

	go func() {
		drop := 0
		for err := shippingEventListener.Listen(ctx, message.HandleFunc(hydrateSagaAndDispatch), "shipping"); err != nil; drop++ {
			log.Warnf("listener exited: %s", err)
		}
	}()

	// 1 queue per saga. uses a dispatcher per message type
	go func() {
		drop := 0
		for err := orderQueueListener.Listen(ctx, message.HandleFunc(hydrateSagaAndDispatch), "order"); err != nil; drop++ {
			log.Warnf("listener exited: %s", err)
		}
	}()

	return nil
}

// hydrateSaga is a middleware handler that takes care of providing the saga state before calling the handler.
// can be implemented in a more generic way via some interfaces. not explored further here.
// this is a standard GET call per saga type to hydrate its state before running a handler.
// will not get more complicated than what is below.
// will benefit from generics in go1.18
func hydrateSagaAndDispatch(ctx context.Context, msg *message.Message) message.Handler {
	// extract the sagaID from the incoming message.
	// we already know the message should be handled by this saga.
	// we can also require a sagaType to dispatch messages from a single queue to multiple saga, but the assumption here
	// is that we keep 1 queue per saga type, because a saga maps to an operation.
	var sagaID string
	if id, ok := msg.Message().UserProperties["sagaID"]; !ok {
		return msg.Error(fmt.Errorf("sagaID not found on message user properties"))
	} else {
		sagaID = id.(string)
	}
	orderSaga := New()
	if err := orderSaga.Hydrate(sagaID); err != nil {
		return msg.Error(err)
	}
	return orderSaga.Dispatcher(ctx, orderSaga, msg)
}

// can certainly be made more generic.
// can be a required func on a Saga interface. not explored further in this example
// will benefit from generics in go1.18
func (*Saga) Dispatcher(ctx context.Context, saga *Saga, msg *message.Message) message.Handler {
	paymentVerified := &messages.PaymentVerified{}
	if msg.Type() == paymentVerified.Type() {
		if err := json.Unmarshal(msg.Message().Data, paymentVerified); err != nil {
			paymentVerified.Message = msg
			return saga.HandlePaymentVerified(ctx, paymentVerified)
		}
	}

	orderShipped := &messages.OrderShipped{}
	if msg.Type() == orderShipped.Type() {
		if err := json.Unmarshal(msg.Message().Data, orderShipped); err != nil {
			orderShipped.Message = msg
			return saga.HandleOrderShipped(ctx, orderShipped)
		}
	}

	createOrderCmd := &messages.CreateOrderCommand{}
	if msg.Type() == createOrderCmd.Type() {
		if err := json.Unmarshal(msg.Message().Data, createOrderCmd); err != nil {
			createOrderCmd.Message = msg
			return saga.HandleCreateOrder(ctx, createOrderCmd)
		}
	}

	completeOrderCmd := &messages.CompleteOrderCommand{}
	if msg.Type() == completeOrderCmd.Type() {
		if err := json.Unmarshal(msg.Message().Data, completeOrderCmd); err != nil {
			completeOrderCmd.Message = msg
			return saga.HandleCompleteOrder(ctx, completeOrderCmd)
		}
	}

	return msg.Error(fmt.Errorf("no handler matching message type"))
}
