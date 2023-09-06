package broker

import (
	"TESTShop"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
	"time"
)

type ConfigNATS struct {
	ClusterID        string
	ClientIdConsumer string
	ClientIdProducer string
	ChannelName      string
	NatsUrl          string
	DurableName      string
	Subject          string
}
type Broker struct {
	sc stan.Conn
}

func NewBroker(sc stan.Conn) *Broker {
	return &Broker{sc: sc}
}

func (b *Broker) ReadFromChannel(channelName string, orders chan<- TESTShop.OrderResponse) error {

	subscription, err := b.sc.Subscribe(channelName, func(msg *stan.Msg) {

		var order TESTShop.OrderResponse
		if err := json.Unmarshal(msg.Data, &order); err != nil {
			logrus.Errorf("error during unmarshall message %s", err)
			return
		}

		orders <- order

	})

	if err != nil {
		logrus.Errorf("error creating subscriber %s", err)
	}

	// Ждем 10 секунд для демонстрации работы
	time.Sleep(10 * time.Second)

	// Отмена подписки и завершение работы
	subscription.Unsubscribe()

	return nil
}
