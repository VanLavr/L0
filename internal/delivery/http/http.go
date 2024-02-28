package http

import "github.com/VanLavr/L0/model"

type HttpHandler struct{}

type Service interface {
	SaveOrder(*model.Order) (string, error)
	GetOrder(string) (*model.Order, error)
	GetOrderIds() []string
}
