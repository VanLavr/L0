package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/VanLavr/L0/internal/delivery/nats"
	"github.com/VanLavr/L0/internal/pkg/config"
	"github.com/VanLavr/L0/internal/pkg/logging"
	"github.com/VanLavr/L0/internal/repo"
	"github.com/VanLavr/L0/internal/service"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	cfg := config.New()
	db := repo.NewPostgres(cfg)
	db.Connect()

	c := repo.NewCache(time.Duration(cfg.Ttl), time.Duration(cfg.Eviction))

	sv := service.New(db, c)
	h := nats.New(sv)

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
