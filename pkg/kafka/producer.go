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
	metadata Metadata
}

func NewProducer(c config.Config) (Producer, error) {
	m, err := ParseMetadata(c)
	if err != nil {
		return nil, err
	}

	cli, err := GetClient(m.Brokers)
	if err != nil {
		return nil, err
	}

	p, err := getProducer(cli)
	if err != nil {
		return nil, fmt.Errorf("error creating producer: %s", err)
	}

	klog.Infof("ready to produce message into %s", m.Topic)

	return &producer{
		producer: p,
		metadata: m,
	}, nil
}

func (p *producer) SendMessage(key string, m interface{}) error {
	value, _ := json.Marshal(&m)

	_, _, err := p.producer.SendMessage(&sarama.ProducerMessage{
		Key:   sarama.StringEncoder(key),
		Topic: p.metadata.Topic,
		Value: sarama.StringEncoder(value),
	})

	return err
}

func getProducer(c sarama.Client) (sarama.SyncProducer, error) {
	c.Config().Producer.RequiredAcks = sarama.WaitForAll
	c.Config().Producer.Retry.Max = 10
	c.Config().Producer.Return.Successes = true
	c.Config().Producer.Idempotent = true
	c.Config().Producer.Return.Errors = true

	producer, err := sarama.NewSyncProducerFromClient(c)
	if err != nil {
		return nil, fmt.Errorf("error creating kafka producer: %s", err)
	}

	return producer, nil
}
