package service

import (
	"TESTShop"
	"TESTShop/pkg/broker"
	"TESTShop/pkg/repository"
)

type Order interface {
	InsertOrderResponse(orderResponse TESTShop.OrderResponse) error
	GetOrderByUID(orderUID string) (*TESTShop.OrderResponse, error)
	ReadFromChannel(cfg broker.ConfigNATSConsumer) error
}

type Service struct {
	Order
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Order: NewOrderService(repos.Order),
	}
}
