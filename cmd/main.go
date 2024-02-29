package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/VanLavr/L0/internal/delivery/http"
	"github.com/VanLavr/L0/internal/delivery/nats"
	"github.com/VanLavr/L0/internal/pkg/config"
	"github.com/VanLavr/L0/internal/pkg/logging"
	"github.com/VanLavr/L0/internal/repo"
	"github.com/VanLavr/L0/internal/service"
	"github.com/labstack/echo/v4"
)

// implement a graceful shutdown (release all resources?)
// make a config
// than initialize repository and database connection
// than initialize a cache
// than initialize a service layer
// init a logger and connection and subsctiption on nats channel
// initialize an echo instance and than initialize http server
// than start the http server
func main() {
	// gs implementation
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// making an new configuration
	cfg := config.New()

	// initializing a repository layer and establishing a database connection
	db := repo.NewPostgres(cfg)
	db.Connect()

	// initializing cache
	c := repo.NewCache(time.Duration(cfg.Ttl), time.Duration(cfg.Eviction))

	// init service
	sv := service.New(db, c, cfg)

	// init nats handler
	h := nats.New(sv)

	logger := logging.New(cfg)
	logger.SetAsDefault()

	// connect and subscribe on nats-streaming
	h.Connect(cfg)
	h.Subscribe(cfg)
	slog.Info("listening channel in nats...")

	// init an echo instance and http server. Run it
	e := echo.New()
	server := http.New(sv)
	http.RegisterRoutes(e, server)
	go func() {
		if err := e.Start(cfg.Addr + ":" + cfg.Port); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}()

	// releasing resources
	<-ctx.Done()
	slog.Info("shutting down")
	h.Unsubscribe()
	h.CloseConnection()
}
