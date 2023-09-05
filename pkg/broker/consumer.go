package broker

import (
	"TESTShop"
	"TESTShop/pkg/cache"
	"TESTShop/pkg/service"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
	"time"
)

type ConfigNATSConsumer struct {
	ClusterID   string
	ClientID    string
	ChannelName string
	NatsUrl     string
	DurableName string
}
type Broker struct {
	Services *service.Service
}

func NewBroker(services *service.Service) *Broker {
	return &Broker{Services: services}
}

func (c *Broker) ReadFromChannel(cfg ConfigNATSConsumer) error {

	sc, err := stan.Connect(cfg.ClusterID, cfg.ClientID)
	if err != nil {
		logrus.Errorf("error to connect nats-streaming server %s", err)
	}
	defer sc.Close()

	subscription, err := sc.Subscribe(cfg.ChannelName, func(msg *stan.Msg) {

		var order TESTShop.OrderResponse
		if err := json.Unmarshal(msg.Data, &order); err != nil {
			logrus.Errorf("error during unmarshall message %s", err)
			return
		}
		cache.AddToCache(&order)
		err = c.Services.InsertOrderResponse(order)
		if err != nil {
			logrus.Errorf("error while writing to database %s", err)
			return
		}
	}, stan.DurableName(cfg.DurableName))

	if err != nil {
		logrus.Errorf("error creating subscriber %s", err)
	}

	// Ждем 10 секунд для демонстрации работы
	time.Sleep(10 * time.Second)

	// Отмена подписки и завершение работы
	subscription.Unsubscribe()
	return nil
}
