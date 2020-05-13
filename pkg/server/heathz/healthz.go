package heathz

import (
	"bennu.cl/identifier-producer/config"
	"bennu.cl/identifier-producer/pkg/kafka"
	"fmt"
	"github.com/Shopify/sarama"
)

type Healthz interface {
	AvailablePartitions() error
	AvailableCluster() error
}

type Admin struct {
	admin    sarama.ClusterAdmin
	client   sarama.Client
	metadata kafka.Metadata
}

func NewHealthz(c config.Config) (Healthz, error) {
	m, err := kafka.ParseMetadata(c)
	if err != nil {
		return nil, fmt.Errorf("error parsing kafka metadata: %s", err)
	}

	cl, err := kafka.GetClient(m.Brokers)
	if err != nil {
		return nil, err
	}

	admin, err := sarama.NewClusterAdminFromClient(cl)
	if err != nil {
		return nil, fmt.Errorf("error creating kafka admin: %s", err)
	}

	return &Admin{
		metadata: m,
		client:   cl,
		admin:    admin,
	}, nil
}

func (a *Admin) AvailablePartitions() error {
	topicsMetadata, err := a.admin.DescribeTopics([]string{a.metadata.Topic})
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

func (a *Admin) AvailableCluster() error {
	brokers := a.metadata.Brokers
	if len(brokers) == 0 {
		return fmt.Errorf("expected at least 1 broker, got 0")
	}

	if err := a.client.RefreshMetadata(a.metadata.Topic); err != nil {
		return fmt.Errorf("error refreshing metadata: %s", err)
	}

	return nil
}
