package service

import (
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"strconv"
	"time"

	"github.com/VanLavr/L0/internal/delivery/http"
	"github.com/VanLavr/L0/internal/pkg/config"
	"github.com/VanLavr/L0/internal/pkg/err"
	"github.com/VanLavr/L0/model"
)

// this package stands for representing business logic of the hole application.
// it has all the features implemented and explained
type service struct {
	repo  Repository
	cache Cache
	cfg   *config.Config
}

type Repository interface {
	// generating id for order
	GenerateTrackNumber() string
	SaveOrder(*model.Order)
	GetOrder(string) (*model.Order, error)
	// connect to postgresql
	Connect() error
	CloseConnection()
	GetIDs() []string
}

type Cache interface {
	Set(key string, value *model.Order, duration time.Duration) error
	Get(key string) (*model.Order, error)
	GetAll() []*model.Order
	Delete(key string) error
}

func New(repo Repository, cache Cache, cfg *config.Config) http.Service {
	return &service{repo: repo, cache: cache, cfg: cfg}
}

// generate id for order, validate it and than save it to database and cache
func (s *service) SaveOrder(order *model.Order) (string, error) {
	// generate id
	uuid := s.repo.GenerateTrackNumber()
	order.Order_uid = uuid

	// validate
	if err := s.validate(order); err != nil {
		slog.Error(err.Error())
		return "", err
	}

	// save
	s.repo.SaveOrder(order)
	if err := s.cache.Set(uuid, order, time.Second*time.Duration(s.cfg.Ttl)); err != nil {
		slog.Error(err.Error())
		return "", err
	}
	return uuid, nil
}

// check if any of the order's fields are empty
func (s *service) validate(order *model.Order) error {
	// order cannot exist without items
	if len(order.Items) == 0 {
		slog.Error(err.ErrEmptyItems.Error())
		return err.ErrEmptyItems
	}

	// get the type of order object to iterate throw all the fields
	tp := reflect.TypeOf(*order)
	// get the value of order
	value := reflect.ValueOf(*order)

	// iterating...
	for i := 0; i < tp.NumField(); i++ {
		// getting the i-ish field and validating it if it's zero value
		if value.Field(i).IsZero() {
			slog.Debug(fmt.Sprintf("field %v is empty", value.Field(i)))
			return err.ErrInvalidField
		}
	}

	return nil
}

// go and fetch data from cache
// if there is such data in cache => respond with this data
// if there no such data in cache, so
// go to database and get it
// than write it to the cache and than => respond with this data
func (s *service) GetOrder(id string) (*model.Order, error) {
	// fetch data from cache
	data, er := s.cache.Get(id)
	if er != nil {
		if errors.Is(er, err.ErrNotFound) {
			// go to database and take data from here
			order, er := s.repo.GetOrder(id)
			if er != nil {
				slog.Error(er.Error())
				return nil, er
			}

			// write data to the cache
			if er = s.cache.Set(order.Order_uid, order, time.Second*time.Duration(s.cfg.Ttl)); er != nil {
				slog.Error(er.Error())
				return nil, er
			}

			// return data from database
			slog.Info("got data from database")
			return order, nil
		} else {
			slog.Error(er.Error())
			return nil, er
		}
	}

	// return data fetched from cache
	slog.Info("got data from cache")
	return data, nil
}

// get the list of all available ids in database
func (s *service) GetOrderIds() []string {
	return s.repo.GetIDs()
}

// get the ids from database
// get all the orders from database with ids
// set orders to the cache
func (s *service) RecoverCache() error {
	// get the ids from database
	ids := s.GetOrderIds()

	// get all orders from database
	// and the GetOrder function will automatically recover cache and set all the orders to
	// the cache, because of write-around strategy
	for _, id := range ids {
		_, err := s.GetOrder(id)
		if err != nil {
			slog.Error(err.Error())
			return err
		}
	}

	slog.Debug(strconv.Itoa(len(s.cache.GetAll())))

	return nil
}
