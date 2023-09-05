package repository

import (
	"TESTShop"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type OrderPostgres struct {
	db *sqlx.DB
}

func NewOrderPostgres(db *sqlx.DB) *OrderPostgres {
	return &OrderPostgres{db: db}
}

func (r *OrderPostgres) InsertOrderResponse(orderResponse TESTShop.OrderResponse) error {
	// Начать транзакцию
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			// Если произошла ошибка, откатываем транзакцию
			logrus.Errorf("error with transaction %s", err)
			tx.Rollback()
			return
		}
		// Если всё успешно, фиксируем транзакцию
		tx.Commit()
	}()

	// Вставляем данные в таблицу delivery_info и получаем ID
	var deliveryID int
	err = tx.QueryRow(`
        INSERT INTO delivery_info (name, phone, zip, city, address, region, email)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id
    `, orderResponse.Delivery.Name, orderResponse.Delivery.Phone, orderResponse.Delivery.Zip, orderResponse.Delivery.City, orderResponse.Delivery.Address, orderResponse.Delivery.Region, orderResponse.Delivery.Email).Scan(&deliveryID)
	if err != nil {
		logrus.Errorf("error while inserting into table delivery_info %s", err)
		return err
	}

	// Вставляем данные в таблицу payment_info и получаем ID
	var paymentID int
	err = tx.QueryRow(`
        INSERT INTO payment_info (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING id
    `, orderResponse.Payment.Transaction, orderResponse.Payment.RequestID, orderResponse.Payment.Currency, orderResponse.Payment.Provider, orderResponse.Payment.Amount, orderResponse.Payment.PaymentDT, orderResponse.Payment.Bank, orderResponse.Payment.DeliveryCost, orderResponse.Payment.GoodsTotal, orderResponse.Payment.CustomFee).Scan(&paymentID)
	if err != nil {
		logrus.Errorf("error while inserting into table payment_info %s", err)
		return err
	}

	// Вставляем данные в таблицу items и получаем ID
	itemIDs := make([]int, len(orderResponse.Items))
	for i, item := range orderResponse.Items {
		err = tx.QueryRow(`
            INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
            RETURNING id
        `, item.ChrtID, item.TrackNumber, item.Price, item.RID, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status).Scan(&itemIDs[i])
		if err != nil {
			logrus.Errorf("error while inserting into table items %s, %d", err, i)
			return err
		}
	}

	// Вставляем данные в таблицу orders
	_, err = tx.Exec(`
        INSERT INTO orders (order_uid, track_number, entry, delivery_id, payment_id, locale, customer_id, delivery_service, shard_key, sm_id, date_created, oof_shard)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
    `, orderResponse.OrderUID, orderResponse.TrackNumber, orderResponse.Entry, deliveryID, paymentID, orderResponse.Locale, orderResponse.CustomerID, orderResponse.DeliveryService, orderResponse.ShardKey, orderResponse.SMID, orderResponse.DateCreated, orderResponse.OofShard)
	if err != nil {
		logrus.Errorf("error while inserting into table orders %s", err)
		return err
	}

	// Вставляем данные в таблицу order_items для связи заказа с товарами
	for _, itemID := range itemIDs {
		_, err = tx.Exec(`
            INSERT INTO order_items (order_id, item_id)
            VALUES ((SELECT id FROM orders WHERE order_uid = $1), $2)
        `, orderResponse.OrderUID, itemID)
		if err != nil {
			logrus.Errorf("error while inserting into table order_items %s", err)
			return err
		}
	}

	return nil
}
func (r *OrderPostgres) ReturnToCache(bool) error {

}
