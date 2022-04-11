package ordersaga

import (
	"context"
	"time"

	"github.com/Azure/go-shuttle/common"
	"github.com/Azure/go-shuttle/message"

	"saga-example/bus"
)

var listeners []common.Listener

type State struct {
	StartedAt         *time.Time
	OrderID           string
	IsOrderShipped    bool
	IsPaymentVerified bool
	IsOrderCompleted  bool
}

type Store interface {
	Save(interface{}) error
	GetOrderSaga(id string) (*State, error)
}

type Saga struct {
	State    *State
	registry map[string]func(ctx context.Context, message interface{}) message.Handler
	bus      bus.Bus
	store    Store
}

func New() *Saga {
	return &Saga{}
}

// Hydrate should be called once, before the message handlers are called.
func (s *Saga) Hydrate(id string) error {
	state, err := s.store.GetOrderSaga(id)
	if err != nil {
		return err
	}
	s.State = state
	return nil
}

func (s *Saga) SaveState() error {
	return s.store.Save(s.State)
}

func (s *Saga) IsOrderSagaCompleted() bool {
	return s.State.IsPaymentVerified && s.State.IsOrderShipped
}
