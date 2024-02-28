package nats

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/VanLavr/L0/internal/delivery/http"
	"github.com/VanLavr/L0/internal/pkg/config"
	"github.com/VanLavr/L0/model"
	"github.com/go-playground/validator/v10"
	"github.com/nats-io/stan.go"
)

// this package stands for connection and interaction with nats-streaming server
// it provides connection and handling nats messages
type Handler struct {
	sc   stan.Conn
	sub  stan.Subscription
	srvc http.Service
}

func New(sv http.Service) *Handler {
	return &Handler{srvc: sv}
}

func (h *Handler) Connect(cfg *config.Config) {
	sc, err := stan.Connect(cfg.ClusterID, cfg.SubName, stan.NatsURL(cfg.NatsUrl))
	if err != nil {
		slog.Debug(err.Error())
		os.Exit(1)
	}

	h.sc = sc
}

func (h *Handler) CloseConnection() {
	if err := h.sc.Close(); err != nil {
		slog.Debug(err.Error())
		os.Exit(1)
	}
}

func (h *Handler) Subscribe(cfg *config.Config) {
	sub, err := h.sc.Subscribe(cfg.SubjName, h.handleMessage)
	if err != nil {
		slog.Debug(err.Error())
		os.Exit(1)
	}
	h.sub = sub
}

func (h *Handler) Unsubscribe() {
	if err := h.sub.Unsubscribe(); err != nil {
		slog.Debug(err.Error())
		os.Exit(1)
	}
}

// unmarshall message to data model
// than validate it's fields with validator
// than save order in database and put it in cache
func (h *Handler) handleMessage(m *stan.Msg) {
	var order model.Order

	// unmarshall
	err := json.Unmarshal(m.Data, &order)
	if err != nil {
		slog.Warn(err.Error())
		return
	}

	// validate
	err = validator.New().Struct(order)
	if err != nil {
		slog.Info("validation failed")
		return
	}

	// save
	id, err := h.srvc.SaveOrder(&order)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	slog.Info("order stored with id: " + id)
}
