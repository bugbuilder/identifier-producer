package kafka

import (
	"bennu.cl/identifier-producer/config"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"k8s.io/klog"
)

type Kafka interface {
	AvailablePartitions() error
	AvailableCluster() error
	GetProducer() (sarama.SyncProducer, error)
	Close() error
}

type kafka struct {
	admin    sarama.ClusterAdmin
	client   sarama.Client
	metadata Metadata
}

var kcl *kafka

func NewKafka(c config.Config) (Kafka, error) {
	if kcl == nil {
		m, err := parseMetadata(c)
		if err != nil {
			return nil, fmt.Errorf("error parsing kafka metadata: %s", err)
		}

		config := sarama.NewConfig()
		config.Version = sarama.V2_4_0_0
		config.Net.MaxOpenRequests = 1

		client, err := sarama.NewClient(m.Brokers, config)
		if err != nil {
			return nil, err
		}
		klog.Infof("connected to %s", m.Brokers)

		admin, err := sarama.NewClusterAdminFromClient(client)
		if err != nil {
			return nil, fmt.Errorf("error creating kafka admin: %s", err)
		}
		kcl = &kafka{
			client:   client,
			admin:    admin,
			metadata: m,
		}
	}
	return kcl, nil
}

func (k *kafka) AvailablePartitions() error {
	topicsMetadata, err := k.admin.DescribeTopics([]string{k.metadata.Topic})
	if err != nil {
		return fmt.Errorf("error describing topics: %s", err)
	}

	if len(topicsMetadata) != 1 {
		return fmt.Errorf("expected only 1 topic metadata, got %d", len(topicsMetadata))
	}

	if len(topicsMetadata[0].Partitions) == 0 {
		return fmt.Errorf("expected at least 1 partition, got 0")
	}

	return nil
}

func (k *kafka) AvailableCluster() error {
	brokers := k.metadata.Brokers
	if len(brokers) == 0 {
		return fmt.Errorf("expected at least 1 broker, got 0")
	}

	if err := k.client.RefreshMetadata(k.metadata.Topic); err != nil {
		return fmt.Errorf("error refreshing metadata: %s", err)
	}

	return nil
}

func (k *kafka) Close() error {
	if err := k.client.Close(); err != nil {
		klog.Infof("failed to shutting down due to client error: %s", err)
		return err
	}
	return nil
}

func (k *kafka) GetProducer() (sarama.SyncProducer, error) {
	k.client.Config().Producer.RequiredAcks = sarama.WaitForAll
	k.client.Config().Producer.Retry.Max = 10
	k.client.Config().Producer.Return.Successes = true
	k.client.Config().Producer.Idempotent = true
	k.client.Config().Producer.Return.Errors = true

	producer, err := sarama.NewSyncProducerFromClient(k.client)
	if err != nil {
		return nil, fmt.Errorf("error creating kafka producer: %s", err)
	}

	return producer, nil
}

func parseMetadata(c config.Config) (Metadata, error) {
	meta := Metadata{}

	if c.Brokers == "" {
		return meta, errors.New("no brokers given")
	}

	if c.Topic == "" {
		return meta, errors.New("no topic given")
	}

	meta.Brokers = []string{c.Brokers}
	meta.Topic = c.Topic

	return meta, nil
}
