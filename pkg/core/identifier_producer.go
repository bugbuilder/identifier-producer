package core

import (
	"bennu.cl/identifier-producer/api/types"
	"bennu.cl/identifier-producer/config"
	"bennu.cl/identifier-producer/pkg/kafka"
	"fmt"
	"k8s.io/klog"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type Message struct {
	types.Identifier
	Node      string `json:"node"`
	Pod       string `json:"pod"`
	CreatedOn string `json:"createdOn"`
}

type identifierProducer struct {
	producer kafka.Producer
}

func NewIdentifierProducer(c config.Config) (Service, error) {
	p, err := kafka.NewProducer(c)
	if err != nil {
		return nil, err
	}

	return &identifierProducer{
		producer: p,
	}, nil
}

func (i *identifierProducer) Save(id types.Identifier) (string, error) {
	key, m := i.getMessage(id)

	err := i.producer.SendMessage(key, m)
	if err != nil {
		klog.Errorf("%s", err)
		return "nil", err
	}

	return key, nil
}

func (i *identifierProducer) getMessage(id types.Identifier) (string, Message) {
	pod := os.Getenv("POD")

	node := os.Getenv("NODE")
	if node == "" {
		node = "localhost"
	}

	if id.TransactionId == "" {
		id.TransactionId = strconv.Itoa(rand.Int())
	}

	key := fmt.Sprintf("%s-%s-%s", node, pod, id.TransactionId)

	return key, Message{
		Identifier: id,
		Node:       node,
		Pod:        pod,
		CreatedOn:  time.Now().Format(time.RFC3339),
	}
}
