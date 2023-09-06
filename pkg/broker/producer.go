package broker

import (
	"TESTShop"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
)

func CreateProducer(clientId, clusterId, natsUrl string) (stan.Conn, error) {

	_, err := nats.Connect(natsUrl)
	if err != nil {
		logrus.Errorf("Error to connect nats-streaming server %s", err)
		return nil, err
	}

	sc, err := stan.Connect(clusterId, clientId)
	if err != nil {
		logrus.Errorf("error to connect nats-streaming producer channel")
		return nil, err
	}

	return sc, nil
}

func PublishMessage(sc stan.Conn, channel string, n int) error {
	for i := 0; i < n; i++ {
		message := generateRandomMessage()
		err := validateOrderResponse(message)
		if err != nil {
			logrus.Errorf("error invalid message %s", err)
			return err
		}
		messageData, err := json.Marshal(message)
		if err != nil {
			return err
		}

		err = sc.Publish(channel, messageData)
		logrus.Print("send messages")
		if err != nil {
			logrus.Errorf("error sending message %s", err)
			return err
		}
	}

	return nil
}
func validateOrderResponse(order TESTShop.OrderResponse) error {
	// Создаем экземпляр валидатора
	validate := validator.New()

	// Проверяем структуру на соответствие правилам валидации
	err := validate.Struct(order)
	if err != nil {
		return err
	}

	return nil
}
