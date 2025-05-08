package kafka

import "github.com/IBM/sarama"

const (
	ProducerPartitionHash = "hash"
	ProducerPartitionRef  = "reference"
	ProducerPartitionCons = "consistent"

	ConsumerGroupRangeBalanceStrategy      = sarama.RangeBalanceStrategyName
	ConsumerGroupRoundRobinBalanceStrategy = sarama.RoundRobinBalanceStrategyName
	ConsumerGroupStickyBalanceStrategy     = sarama.StickyBalanceStrategyName
)

// kafka config struct
type Config struct {
	Addresses             []string `yaml:"addresses,omitempty"`
	ProducerPartitioner   string   `yaml:"producerPartitioner,omitempty"`
	ConsumerGroupStrategy string   `yaml:"consumerGroupStrategy,omitempty"`
}
