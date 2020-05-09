package api

import (
	"bennu.cl/identifier-producer/config"
	"bennu.cl/identifier-producer/pkg/api"
	"bennu.cl/identifier-producer/pkg/kafka"
	"fmt"
	"k8s.io/klog"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type Message struct {
	api.Identifier
	Node      string `json:"node"`
	Pod       string `json:"pod"`
	CreatedOn string `json:"createdOn"`
}

type identifier struct {
	producer kafka.Producer
}

func NewIdentifierProducer(c config.Config) (api.IdentifierService, error) {
	p, err := kafka.NewProducer(c)
	if err != nil {
		return nil, err
	}

	return &identifier{
		producer: p,
	}, nil
}

func (i *identifier) Save(id api.Identifier) (string, error) {
	key, m := i.getMessage(id)

	err := i.producer.SendMessage(key, m)
	if err != nil {
		klog.Errorf("%s", err)
		return "nil", err
	}

	return key, nil
}

func (i *identifier) getMessage(id api.Identifier) (string, Message) {
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
