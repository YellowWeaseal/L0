package cache

import (
	"TESTShop"
	"github.com/sirupsen/logrus"
	"sync"
)

type Cache struct {
	cache      map[string]TESTShop.OrderResponse
	cacheMutex *sync.Mutex
}

func NewCache(cache map[string]TESTShop.OrderResponse, cacheMutex *sync.Mutex) *Cache {
	//map mutex
	return &Cache{cache: cache, cacheMutex: cacheMutex}
}

func (c *Cache) AddToCache(order TESTShop.OrderResponse) {

	c.cacheMutex.Lock()
	defer c.cacheMutex.Unlock()

	c.cache[order.OrderUID] = order
}

func (c *Cache) GetOrderByUID(orderUID string) (TESTShop.OrderResponse, bool) {

	c.cacheMutex.Lock()
	defer c.cacheMutex.Unlock()

	order, found := c.cache[orderUID]
	if !found {
		logrus.Error("order is not found in cache ")
		return order, false
	}

	return order, true
}
