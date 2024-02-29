package http

import (
	"github.com/VanLavr/L0/model"
	"github.com/labstack/echo/v4"
)

type HttpHandler struct {
	svc Service
}

type Service interface {
	SaveOrder(*model.Order) (string, error)
	GetOrder(string) (*model.Order, error)
	GetOrderIds() []string
	RecoverCache() error
}

func New(svc Service) *HttpHandler {
	return &HttpHandler{svc: svc}
}

func RegisterRoutes(e *echo.Echo, srv *HttpHandler) {
	e.GET("/order/ids", srv.GetIds)
	e.GET("/order", srv.GetOrder)
}

func (h *HttpHandler) GetIds(c echo.Context) error {
	panic("not implemented")
}

func (h *HttpHandler) GetOrder(c echo.Context) error {
	panic("not implemented")
}
