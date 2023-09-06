package service

import (
	"TESTShop"
	"TESTShop/pkg/broker"
	"TESTShop/pkg/repository"
	"fmt"
	"github.com/sirupsen/logrus"
)

type Cache interface {
	GetOrderByUID(orderUID string) (TESTShop.OrderResponse, bool)
	AddToCache(order TESTShop.OrderResponse)
}

type OrderService struct {
	repo   repository.Order
	broker *broker.Broker
	cache  Cache
}

func (s *OrderService) InitCache() {
	//ичиатем из бд все значения ордеров и заполняем в кэш
	orders := make(chan TESTShop.OrderResponse)

	go func() {
		defer close(orders)

		err := s.repo.GetAllOrderResponses(orders)
		if err != nil {
			logrus.Errorf("error during cache initialization %s", err)
		}

	}()

	for orderResponse := range orders {
		s.cache.AddToCache(orderResponse)
	}

	return

}

func NewOrderService(repo repository.Order, broker *broker.Broker, cache Cache) *OrderService {
	return &OrderService{repo: repo, broker: broker, cache: cache}
}

func (s *OrderService) InsertOrderResponse(orderResponse TESTShop.OrderResponse) error {
	err := s.repo.InsertOrderResponse(orderResponse)
	if err != nil {
		return fmt.Errorf("repo.InsertOrderResponse : %w", err)
	}
	return s.repo.InsertOrderResponse(orderResponse)
}
func (s *OrderService) GetOrderByUID(orderUID string) (*TESTShop.OrderResponse, error) {

	order, found := s.cache.GetOrderByUID(orderUID)
	if !found {

		order, err := s.repo.GetOrderByIdPostgres(orderUID)
		if err != nil {
			logrus.Errorf("order not found in database %s", err)
		}
		s.cache.AddToCache(order)

		return nil, nil
	}

	return &order, nil
}
func (s *OrderService) StartProcessMessages(channelName string) error {
	orders := make(chan TESTShop.OrderResponse)

	go func() {
		err := s.broker.ReadFromChannel(channelName, orders)
		if err != nil {
			logrus.Errorf("error while reading from channel %s", err)
		}
	}()

	for order := range orders {
		s.cache.AddToCache(order)
		err := s.repo.InsertOrderResponse(order)
		if err != nil {
			logrus.Errorf("error while writing to database %s", err)
		}
	}
	return nil
}
