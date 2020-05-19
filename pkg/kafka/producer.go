package kafka

import (
	"bennu.cl/identifier-producer/config"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"k8s.io/klog"
)

type Producer interface {
	SendMessage(key string, value interface{}) error
}

type producer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewProducer(c config.Config) (Producer, error) {
	k, err := NewKafka(c)
	if err != nil {
		return nil, err
	}

	p, err := k.GetProducer()
	if err != nil {
		return nil, fmt.Errorf("error creating producer: %s", err)
	}

	klog.Infof("ready to produce message into %s", c.Topic)

	return &producer{
		producer: p,
		topic:    c.Topic,
	}, nil
}

func (p *producer) SendMessage(key string, m interface{}) error {
	value, _ := json.Marshal(&m)

	_, _, err := p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.StringEncoder(value),
	})

	return err
}
