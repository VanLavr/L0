package repo

import (
	"sync"
	"time"

	"github.com/VanLavr/L0/internal/pkg/err"
	"github.com/VanLavr/L0/internal/service"
	"github.com/VanLavr/L0/model"
)

type item struct {
	Order          model.Order
	ExpirationTime int64
	Created        time.Time
}

type cache struct {
	sync.RWMutex
	defaultExpiration time.Duration
	cacheEviction     time.Duration
	items             map[string]item
}

func New(defaultExpiration, cacheEviction time.Duration) service.Cache {
	c := cache{
		defaultExpiration: defaultExpiration,
		cacheEviction:     cacheEviction,
		items:             make(map[string]item),
	}

	if cacheEviction > 0 {
		c.startGC()
	}

	return &c
}

func (c *cache) Set(key string, value model.Order, duration time.Duration) {
	var exp int64

	if duration > 0 {
		exp = time.Now().Add(duration).UnixNano()
	}

	c.Lock()
	defer c.Unlock()

	c.items[key] = item{
		Order:          value,
		ExpirationTime: exp,
		Created:        time.Now(),
	}
}

func (c *cache) Get(key string) (model.Order, error) {
	c.RLock()
	defer c.RUnlock()

	item, ok := c.items[key]
	if !ok {
		return model.Order{}, err.ErrNotFound
	}

	if item.ExpirationTime > 0 {
		if time.Now().UnixNano() > item.ExpirationTime {
			return model.Order{}, err.ErrNotFound
		}
	}

	return item.Order, nil
}

func (c *cache) Delete(key string) error {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.items[key]; !ok {
		return err.ErrNoSuchKeyInCache
	}

	delete(c.items, key)

	return nil
}

func (c *cache) startGC() {
	go c.GC()
}

func (c *cache) GC() {
	for {
		<-time.After(c.cacheEviction)
		if c.items == nil {
			return
		}

		if keys := c.expiredKeys(); len(keys) != 0 {
			c.clearItems(keys)
		}
	}
}

func (c *cache) expiredKeys() []string {
	c.RLock()
	defer c.RUnlock()

	var keys []string
	for k, i := range c.items {
		if time.Now().UnixNano() > i.ExpirationTime && i.ExpirationTime > 0 {
			keys = append(keys, k)
		}
	}

	return keys
}

func (c *cache) clearItems(keys []string) {
	c.Lock()
	defer c.Unlock()

	for _, k := range keys {
		delete(c.items, k)
	}
}
