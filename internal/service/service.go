package service

import (
	"fmt"
	"log/slog"
	"reflect"
	"time"

	"github.com/VanLavr/L0/internal/delivery/http"
	"github.com/VanLavr/L0/internal/pkg/err"
	"github.com/VanLavr/L0/model"
)

type service struct {
	repo Repository
}

type Repository interface {
	GenerateTrackNumber() string
	SaveOrder(*model.Order)
	GetOrder(string) (*model.Order, error)
	Connect() error
}

type Cache interface {
	Set(key string, value *model.Order, duration time.Duration)
	Get(key string) (*model.Order, error)
	Delete(key string) error
}

func New(repo Repository) http.Service {
	return &service{repo: repo}
}

func (s *service) SaveOrder(order *model.Order) (string, error) {
	uuid := s.repo.GenerateTrackNumber()
	order.Order_uid = uuid

	if err := s.validate(order); err != nil {
		return "", err
	}

	s.repo.SaveOrder(order)
	return uuid, nil
}

func (s *service) validate(order *model.Order) error {
	if len(order.Items) == 0 {
		return err.ErrEmptyItems
	}

	tp := reflect.TypeOf(*order)
	value := reflect.ValueOf(*order)

	for i := 0; i < tp.NumField(); i++ {
		if value.Field(i).IsZero() {
			slog.Debug(fmt.Sprintf("field %v is empty", value.Field(i)))
			return err.ErrInvalidField
		}
	}

	return nil
}

func (s *service) GetOrder(string) *model.Order {
	panic("not implemented")
}

func (s *service) GetOrderIds() []string {
	panic("not implemented")
}
