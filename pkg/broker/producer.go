package broker

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
)

type ConfigNATSProducer struct {
	ClusterID string
	ClientID  string
	NatsUrl   string
}

func CreateProducer(cfg ConfigNATSProducer) (stan.Conn, error) {
	nc, err := nats.Connect(cfg.NatsUrl)
	if err != nil {
		logrus.Errorf("Error to connect nats-streaming server %s", err)
		return nil, err
	}

	sc, err := stan.Connect(cfg.ClusterID, cfg.ClientID, stan.NatsConn(nc))
	if err != nil {
		logrus.Errorf("error to connect nats-streaming producer channel")
		return nil, err
	}

	return sc, nil
}

func SendMessage(sc stan.Conn, subject string, n int) error {
	for i := 0; i < n; i++ {
		message := generateRandomMessage()
		messageData, err := json.Marshal(message)
		if err != nil {
			return err
		}

		err = sc.Publish(subject, messageData)
		if err != nil {
			logrus.Errorf("error sending message %s", err)
			return err
		}
	}

	return nil
}
