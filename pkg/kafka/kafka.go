package kafka

import (
	"bennu.cl/identifier-producer/config"
	"errors"
	"github.com/Shopify/sarama"
	"k8s.io/klog"
	"os"
	"os/signal"
	"syscall"
)

var cli sarama.Client

func getClient(brokers []string) (sarama.Client, error) {
	if cli == nil {
		config := sarama.NewConfig()
		config.Version = sarama.V2_4_0_0
		config.Net.MaxOpenRequests = 1

		client, err := sarama.NewClient(brokers, config)
		if err != nil {
			return nil, err
		}

		klog.Infof("connected to %s", brokers)
		cli = client

		close(brokers)

		return cli, nil
	}

	return cli, nil
}

func close(b []string) {
	go func() {
		quit := make(chan os.Signal)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		if err := cli.Close(); err != nil {
			klog.Infof("failed to shut down client: %s", err)
		}

		klog.Infof("disconnected from %s", b)
	}()
}

func ParseMetadata(c config.Config) (Metadata, error) {
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
