package kafka

import (
	"github.com/IBM/sarama"
)

func NewClusterAdmin(c Config) (admin sarama.ClusterAdmin, err error) {
	return sarama.NewClusterAdmin(c.Addresses, nil)
}

func NewClient(c Config) (client sarama.Client, err error) {
	return sarama.NewClient(c.Addresses, nil)
}
