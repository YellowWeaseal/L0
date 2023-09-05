package cache

import (
	"TESTShop"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"sync"
)

/*type Cache struct {
	Services *service.Service
}

func NewCache(services *service.Service) *Cache {
	return &Cache{Services: services}
}*/

var cacheMutex sync.Mutex
var cache map[string]*TESTShop.OrderResponse

func AddToCache(order *TESTShop.OrderResponse) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// Сохраняем заказ в кэше по ключу OrderUID
	cache[order.OrderUID] = order
}

func GetOrderByUID(orderUID string) (*TESTShop.OrderResponse, error) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// Извлекаем заказ из кэша по ключу OrderUID
	order, found := cache[orderUID]
	if !found {
		logrus.Error("order is not found in cache ")
		return nil, errors.New("order not found")
	}

	return order, nil
}
