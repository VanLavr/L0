package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/VanLavr/L0/internal/nats"
	"github.com/VanLavr/L0/internal/pkg/config"
	"github.com/VanLavr/L0/internal/pkg/logging"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	h := new(nats.Handler)
	cfg := config.New()
	logger := logging.New(cfg)
	logger.SetAsDefault()

	h.Connect(cfg)
	h.Subscribe(cfg)
	slog.Info("listening channel in nats...")

	<-ctx.Done()
	slog.Info("shutting down")
	h.Unsubscribe()
	h.CloseConnection()
}
