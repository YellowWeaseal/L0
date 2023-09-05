package repository

import (
	"TESTShop"
	"github.com/jmoiron/sqlx"
)

type Order interface {
	InsertOrderResponse(orderResponse TESTShop.OrderResponse) error
	ReturnToCache(bool) error
}

type Repository struct {
	Order
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Order: NewOrderPostgres(db),
	}
}
