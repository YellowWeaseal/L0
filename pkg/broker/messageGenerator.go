package broker

import (
	"TESTShop"
	"fmt"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func generateRandomMessage() TESTShop.OrderResponse {

	rand.Seed(time.Now().UnixNano())

	message := TESTShop.OrderResponse{
		OrderUID:    randomString(16),
		TrackNumber: fmt.Sprintf("WBILMTESTTRACK-%d", rand.Intn(10000)),
		Entry:       "WBIL",
		Delivery: TESTShop.Delivery{
			Name:    randomString(12),
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   randomString(10) + "@gmail.com",
		},
		Payment: TESTShop.Payment{
			Transaction:  randomString(16),
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       rand.Intn(5000) + 1000,
			PaymentDT:    time.Now().Unix(),
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   rand.Intn(500),
			CustomFee:    0,
		},
		Items: []TESTShop.Item{
			{
				ChrtID:      rand.Intn(1000000),
				TrackNumber: fmt.Sprintf("WBILMTESTTRACK-%d", rand.Intn(10000)),
				Price:       rand.Intn(500),
				RID:         randomString(16),
				Name:        "Mascaras",
				Sale:        rand.Intn(50),
				Size:        "0",
				TotalPrice:  rand.Intn(500),
				NmID:        rand.Intn(1000000),
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
		},
		Locale:            "en",
		InternalSignature: randomString(20),
		CustomerID:        randomString(8),
		DeliveryService:   "meest",
		ShardKey:          fmt.Sprintf("%d", rand.Intn(10)),
		SMID:              rand.Intn(100),
		DateCreated:       time.Now(),
		OofShard:          randomString(1),
	}
	logrus.Printf("generate order %s", message)
	return message
}
