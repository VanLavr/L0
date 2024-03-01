package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
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
// init a logger
// than initialize repository and database connection
// than initialize a cache and do cache recover
// than initialize a service layer
// init connection subsctiption on nats channel
// initialize an echo instance and than initialize http server
// than start the http server
func main() {
	// gs implementation
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// making an new configuration
	cfg := config.New()

	// init a logger
	logger := logging.New(cfg)
	logger.SetAsDefault()

	// initializing a repository layer and establishing a database connection
	db := repo.NewPostgres(cfg)
	db.Connect()

	// initializing cache
	c := repo.NewCache(time.Duration(cfg.Ttl), time.Duration(cfg.Eviction))
	slog.Debug(strconv.Itoa(len(c.GetAll())))
	for _, order := range c.GetAll() {
		slog.Debug(order.Order_uid)
	}

	// init service
	sv := service.New(db, c, cfg)
	if err := sv.RecoverCache(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	// init nats handler
	h := nats.New(sv)

	// connect and subscribe on nats-streaming
	h.Connect(cfg)
	h.Subscribe(cfg)
	slog.Info("listening channel in nats...")

	// init an echo instance and http server. Run it
	e := echo.New()
	e.Static("/static", "static")
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
	db.CloseConnection()
}
