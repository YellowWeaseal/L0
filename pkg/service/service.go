package service

import (
	"TESTShop"
)

type Order interface {
	InsertOrderResponse(orderResponse TESTShop.OrderResponse) error
	GetOrderByUID(orderUID string) (*TESTShop.OrderResponse, error)
	StartProcessMessages(channelName, durableName string) error
}
