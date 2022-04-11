package bus

import "golang.org/x/net/context"

type Bus interface {
	PublishEvent(ctx context.Context, event interface{}) error
	SendCommand(ctx context.Context, command interface{}) error
}
