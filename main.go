package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/prometheus/common/log"

	"saga-example/ordersaga"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)

	if err := ordersaga.StartListeners(ctx); err != nil {
		log.Fatalf("failed to setup listeners: %s", err)
	}

	// start other sagas/listeners here.

	select {
	case <-ctx.Done():
		log.Info(ctx.Err())
		stop()
	}
}
