package service

import (
	"TESTShop"
	"TESTShop/pkg/broker"
	"TESTShop/pkg/repository"
)

type OrderService struct {
	repo repository.Order
}

func NewOrderService(repo repository.Order) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) InsertOrderResponse(orderResponse TESTShop.OrderResponse) error {
	return s.repo.InsertOrderResponse(orderResponse)
}
func (s *OrderService) GetOrderByUID(orderUID string) (*TESTShop.OrderResponse, error) {
	return s.GetOrderByUID(orderUID)
}
func (s *OrderService) ReadFromChannel(cfg broker.ConfigNATSConsumer) error {
	return s.ReadFromChannel(cfg)
}
