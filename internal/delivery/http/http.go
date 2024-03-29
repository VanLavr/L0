package http

import (
	"errors"
	"fmt"
	"log/slog"

	er "github.com/VanLavr/L0/internal/pkg/err"
	"github.com/VanLavr/L0/model"
	"github.com/VanLavr/L0/view/errorview"
	IDS "github.com/VanLavr/L0/view/ids"
	"github.com/VanLavr/L0/view/orders"
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

// just return the ids
func (h *HttpHandler) GetIds(c echo.Context) error {
	ids := h.svc.GetOrderIds()
	var params []string
	for _, id := range ids {
		param := fmt.Sprintf("/order?order_uid=%s", id)
		params = append(params, param)
	}
	return Render(c, IDS.ShowIds(ids, params))
}

// parse query params
// than find the order with specified id
// than return it
func (h *HttpHandler) GetOrder(c echo.Context) error {
	// pars query params
	id := c.QueryParam("order_uid")
	slog.Debug(id)

	// find the order with specified id
	order, err := h.svc.GetOrder(id)
	if err != nil {
		if errors.Is(err, er.ErrNotFound) {
			return Render(c, errorview.ShowError(err))
		} else {
			return Render(c, errorview.ShowError(err))
		}
	}

	// return order
	return Render(c, orders.ShowOrder(*order))
}
