package service

import (
	"time"

	"github.com/VanLavr/L0/internal/delivery/http"
	"github.com/VanLavr/L0/model"
)

type service struct{}

type Repository interface {
	SaveOrder(*model.Order) string
	GetOrder(string) *model.Order
}

type Cache interface {
	Set(key string, value model.Order, duration time.Duration)
	Get(key string) (model.Order, error)
	Delete(key string) error
}

func New() http.Service {
	return &service{}
}

func (s *service) SaveOrder(*model.Order) (string, error) {
	panic("not implemented")
}

func (s *service) GetOrder(string) *model.Order {
	panic("not implemented")
}

func (s *service) GetOrderIds() []string {
	panic("not implemented")
}
