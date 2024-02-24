package nats

import (
	"encoding/json"
	"log"

	"github.com/VanLavr/L0/model"
	"github.com/go-playground/validator/v10"
	"github.com/nats-io/stan.go"
)

type Handler struct {
	Name      string
	ClusterID string
}

func (h *Handler) HandleMessage(m *stan.Msg) model.Order {
	var order model.Order
	err := json.Unmarshal(m.Data, &order)
	if err != nil {
		log.Fatal(err)
	}

	err = validator.New().Struct(order)
	if err != nil {
		log.Fatal(err)
	}

	return order
}
