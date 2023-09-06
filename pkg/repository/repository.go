package repository

import (
	"TESTShop"
	"github.com/jmoiron/sqlx"
)

type Order interface {
	InsertOrderResponse(orderResponse TESTShop.OrderResponse) error
	GetOrderByIdPostgres(orderUID string) (TESTShop.OrderResponse, error)
	GetAllOrderResponses(orders chan<- TESTShop.OrderResponse) error
}

type Repository struct {
	Order
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Order: NewOrderPostgres(db),
	}
}
