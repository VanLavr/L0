package service

import (
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"time"

	"github.com/VanLavr/L0/internal/delivery/http"
	"github.com/VanLavr/L0/internal/pkg/config"
	"github.com/VanLavr/L0/internal/pkg/err"
	"github.com/VanLavr/L0/model"
)

type service struct {
	repo  Repository
	cache Cache
	cfg   *config.Config
}

type Repository interface {
	GenerateTrackNumber() string
	SaveOrder(*model.Order)
	GetOrder(string) (*model.Order, error)
	Connect() error
	GetIDs() []string
}

type Cache interface {
	Set(key string, value *model.Order, duration time.Duration)
	Get(key string) (*model.Order, error)
	Delete(key string) error
}

func New(repo Repository, cache Cache) http.Service {
	return &service{repo: repo, cache: cache}
}

func (s *service) SaveOrder(order *model.Order) (string, error) {
	uuid := s.repo.GenerateTrackNumber()
	order.Order_uid = uuid

	if err := s.validate(order); err != nil {
		return "", err
	}

	s.repo.SaveOrder(order)
	s.cache.Set(uuid, order, time.Second*time.Duration(s.cfg.Ttl))
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

// go and fetch data from cache
// if there is such data in database => respond with this data
// if there no such data in cache, so
// go to database and get it
// than write it to the cache and than => respond with this data
func (s *service) GetOrder(id string) *model.Order {
	order, err := s.repo.GetOrder(id)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	return order
}

func (s *service) GetOrderIds() []string {
	return s.repo.GetIDs()
}
