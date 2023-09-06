package repository

import (
	"TESTShop"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"log"
)

type OrderPostgres struct {
	db *sqlx.DB
}

func NewOrderPostgres(db *sqlx.DB) *OrderPostgres {
	return &OrderPostgres{db: db}
}

func (r *OrderPostgres) InsertOrderResponse(orderResponse TESTShop.OrderResponse) error {

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			logrus.Errorf("error with transaction %s", err)
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

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

	_, err = tx.Exec(`
        INSERT INTO orders (order_uid, track_number, entry, delivery_id, payment_id, locale, internal_signature,customer_id, delivery_service, shard_key, sm_id, date_created, oof_shard)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
    `, orderResponse.OrderUID, orderResponse.TrackNumber, orderResponse.Entry, deliveryID, paymentID, orderResponse.Locale, orderResponse.InternalSignature, orderResponse.CustomerID, orderResponse.DeliveryService, orderResponse.ShardKey, orderResponse.SMID, orderResponse.DateCreated, orderResponse.OofShard)
	if err != nil {
		logrus.Errorf("error while inserting into table orders %s", err)
		return fmt.Errorf("error while inserting into table orders %w", err)
	}

	for _, itemID := range itemIDs {
		_, err = tx.Exec(`
            INSERT INTO order_items (order_id, item_id)
            VALUES ($1, $2)
        `, orderResponse.OrderUID, itemID)
		if err != nil {
			logrus.Errorf("error while inserting into table order_items %s", err)
			return err
		}
	}

	return nil
}
func (r *OrderPostgres) GetOrderByIdPostgres(orderUID string) (TESTShop.OrderResponse, error) {

	queryOrder := `
		SELECT
			o.order_uid,
			o.track_number,
			o.entry,
			o.locale,
			o.internal_signature,
			o.customer_id,
			o.delivery_service,
			o.shard_key,
			o.sm_id,
			o.date_created,
			o.oof_shard,
			d.name AS delivery_name,
			d.phone AS delivery_phone,
			d.zip AS delivery_zip,
			d.city AS delivery_city,
			d.address AS delivery_address,
			d.region AS delivery_region,
			d.email AS delivery_email,
			p.transaction,
			p.request_id,
			p.currency,
			p.provider,
			p.amount,
			p.payment_dt,
			p.bank,
			p.delivery_cost,
			p.goods_total,
			p.custom_fee
		FROM
			orders o
		INNER JOIN
			delivery_info d ON o.delivery_id = d.id
		INNER JOIN
			payment_info p ON o.payment_id = p.id
		WHERE
			o.order_uid = $1
	`

	queryItems := `
		SELECT
			i.chrt_id,
			i.track_number,
			i.price,
			i.rid,
			i.name,
			i.sale,
			i.size,
			i.total_price,
			i.nm_id,
			i.brand,
			i.status
		FROM
			order_items oi
		INNER JOIN
			items i ON oi.item_id = i.id
		WHERE
			oi.order_id = $1
	`

	// Замените на фактический order_uid

	var orderResponse TESTShop.OrderResponse
	err := r.db.QueryRow(queryOrder, orderUID).Scan(
		&orderResponse.OrderUID,
		&orderResponse.TrackNumber,
		&orderResponse.Entry,
		&orderResponse.Locale,
		&orderResponse.InternalSignature,
		&orderResponse.CustomerID,
		&orderResponse.DeliveryService,
		&orderResponse.ShardKey,
		&orderResponse.SMID,
		&orderResponse.DateCreated,
		&orderResponse.OofShard,
		&orderResponse.Delivery.Name,
		&orderResponse.Delivery.Phone,
		&orderResponse.Delivery.Zip,
		&orderResponse.Delivery.City,
		&orderResponse.Delivery.Address,
		&orderResponse.Delivery.Region,
		&orderResponse.Delivery.Email,
		&orderResponse.Payment.Transaction,
		&orderResponse.Payment.RequestID,
		&orderResponse.Payment.Currency,
		&orderResponse.Payment.Provider,
		&orderResponse.Payment.Amount,
		&orderResponse.Payment.PaymentDT,
		&orderResponse.Payment.Bank,
		&orderResponse.Payment.DeliveryCost,
		&orderResponse.Payment.GoodsTotal,
		&orderResponse.Payment.CustomFee,
	)

	if err != nil {
		log.Fatal(err)
	}

	rows, err := r.db.Query(queryItems, orderUID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var items []TESTShop.Item
	for rows.Next() {
		var item TESTShop.Item
		err := rows.Scan(
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.RID,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, item)
	}

	orderResponse.Items = items

	return orderResponse, nil
}

func (r *OrderPostgres) GetAllOrderResponses(orders chan<- TESTShop.OrderResponse) error {

	// Выполните SQL-запрос для выборки всех заказов и связанных данных
	query := `
		SELECT
			o.order_uid,
			o.track_number,
			o.entry,
			o.locale,
			o.customer_id,
			o.delivery_service,
			o.shard_key,
			o.sm_id,
			o.date_created,
			o.oof_shard,
			d.name AS delivery_name,
			d.phone AS delivery_phone,
			d.zip AS delivery_zip,
			d.city AS delivery_city,
			d.address AS delivery_address,
			d.region AS delivery_region,
			d.email AS delivery_email,
			p.transaction,
			p.request_id,
			p.currency,
			p.provider,
			p.amount,
			p.payment_dt,
			p.bank,
			p.delivery_cost,
			p.goods_total,
			p.custom_fee
		FROM
			orders o
		INNER JOIN
			delivery_info d ON o.delivery_id = d.id
		INNER JOIN
			payment_info p ON o.payment_id = p.id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var orderResponse TESTShop.OrderResponse
		err := rows.Scan(
			&orderResponse.OrderUID,
			&orderResponse.TrackNumber,
			&orderResponse.Entry,
			&orderResponse.Locale,
			&orderResponse.InternalSignature,
			&orderResponse.CustomerID,
			&orderResponse.DeliveryService,
			&orderResponse.ShardKey,
			&orderResponse.SMID,
			&orderResponse.DateCreated,
			&orderResponse.OofShard,
			&orderResponse.Delivery.Name,
			&orderResponse.Delivery.Phone,
			&orderResponse.Delivery.Zip,
			&orderResponse.Delivery.City,
			&orderResponse.Delivery.Address,
			&orderResponse.Delivery.Region,
			&orderResponse.Delivery.Email,
			&orderResponse.Payment.Transaction,
			&orderResponse.Payment.RequestID,
			&orderResponse.Payment.Currency,
			&orderResponse.Payment.Provider,
			&orderResponse.Payment.Amount,
			&orderResponse.Payment.PaymentDT,
			&orderResponse.Payment.Bank,
			&orderResponse.Payment.DeliveryCost,
			&orderResponse.Payment.GoodsTotal,
			&orderResponse.Payment.CustomFee,
		)
		if err != nil {
			log.Println(err)
			continue
		}

		// Добавьте данные о товарах, связанных с заказом
		orderResponse.Items, err = r.GetItemsForOrder(orderResponse.OrderUID)
		if err != nil {
			log.Println(err)
			continue
		}

		// Отправьте OrderResponse в канал
		orders <- orderResponse
	}

	return nil
}

func (r *OrderPostgres) GetItemsForOrder(orderUID string) ([]TESTShop.Item, error) {
	// Выполните SQL-запрос для выборки товаров, связанных с заказом
	queryItems := `
		SELECT
			i.chrt_id,
			i.track_number,
			i.price,
			i.rid,
			i.name,
			i.sale,
			i.size,
			i.total_price,
			i.nm_id,
			i.brand,
			i.status
		FROM
			order_items oi
		INNER JOIN
			items i ON oi.item_id = i.id
		WHERE
			oi.order_id = $1
	`

	rows, err := r.db.Query(queryItems, orderUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []TESTShop.Item
	for rows.Next() {
		var item TESTShop.Item
		err := rows.Scan(
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.RID,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

/*func main() {
	// Установите подключение к базе данных PostgreSQL
	connStr := "user=username dbname=yourdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создайте канал для передачи данных
	ch := make(chan OrderResponse)

	// Запустите отдельную горутину для обработки данных из канала
	go func() {
		for orderResponse := range ch {
			// Здесь вы можете обрабатывать каждый OrderResponse по вашему усмотрению
			fmt.Printf("%+v\n", orderResponse)
		}
	}()

	// Вызов функции для получения всех OrderResponse из базы данных и передачи их в канал
	err = GetAllOrderResponses(db, ch)
	if err != nil {
		log.Fatal(err)
	}

	// Подождите, пока обработка данных завершится
	time.Sleep(time.Second * 2)
}*/
