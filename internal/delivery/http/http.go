package http

import (
	"errors"
	"net/http"

	er "github.com/VanLavr/L0/internal/pkg/err"
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

// just return the ids
func (h *HttpHandler) GetIds(c echo.Context) error {
	ids := h.svc.GetOrderIds()
	return c.JSON(http.StatusOK, Response{
		Error:   "",
		Content: ids,
	})
}

// parse query params
// than find the order with specified id
// than return it
func (h *HttpHandler) GetOrder(c echo.Context) error {
	// pars query params
	id := c.QueryParam("order_uid")

	// find the order with specified id
	order, err := h.svc.GetOrder(id)
	if err != nil {
		if errors.Is(err, er.ErrNotFound) {
			return c.JSON(http.StatusNotFound, Response{
				Error:   err.Error(),
				Content: nil,
			})
		} else {
			return c.JSON(http.StatusInternalServerError, Response{
				Error:   er.ErrInternal.Error(),
				Content: nil,
			})
		}
	}

	// return order
	return c.JSON(http.StatusFound, Response{
		Error:   "",
		Content: order,
	})
}
