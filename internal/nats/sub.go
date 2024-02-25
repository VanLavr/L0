package nats

import (
	"encoding/json"
	"log/slog"

	"github.com/VanLavr/L0/internal/delivery/http"
	"github.com/VanLavr/L0/internal/pkg/config"
	"github.com/VanLavr/L0/model"
	"github.com/go-playground/validator/v10"
	"github.com/nats-io/stan.go"
)

type Handler struct {
	Name      string
	ClusterID string
	NatsURL   string
	sc        stan.Conn
	sub       stan.Subscription
	srvc      http.Service
}

func New(name, clusterid, natsurl string) *Handler {
	return &Handler{Name: name, ClusterID: clusterid, NatsURL: natsurl}
}

func (h *Handler) Connect(cfg *config.Config) {
	sc, err := stan.Connect(cfg.ClusterID, cfg.SubName, stan.NatsURL(cfg.NatsUrl))
	if err != nil {
		slog.Debug(err.Error())
	}

	h.sc = sc
}

func (h *Handler) CloseConnection() {
	if err := h.sc.Close(); err != nil {
		slog.Debug(err.Error())
	}
}

func (h *Handler) Subscribe(cfg *config.Config) {
	sub, err := h.sc.Subscribe(cfg.SubjName, h.handleMessage)
	if err != nil {
		slog.Debug(err.Error())
	}
	h.sub = sub
}

func (h *Handler) Unsubscribe() {
	if err := h.sub.Unsubscribe(); err != nil {
		slog.Debug(err.Error())
	}
}

func (h *Handler) handleMessage(m *stan.Msg) {
	var order model.Order
	err := json.Unmarshal(m.Data, &order)
	if err != nil {
		slog.Warn(err.Error())
		return
	}

	err = validator.New().Struct(order)
	if err != nil {
		slog.Info("validation failed")
		return
	}

	id, err := h.srvc.SaveOrder(&order)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	slog.Info("order stored with id: " + id)
}
