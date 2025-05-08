package kafka

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/IBM/sarama"
)

var c = Config{
	Addresses:             []string{"localhost:9092", "localhost:9093", "localhost:9094"},
	ProducerPartitioner:   "hash",
	ConsumerGroupStrategy: "range",
}

func TestNewProducer(t *testing.T) {
	pd, err := NewProducer(c)
	if err != nil {
		t.Errorf("NewProducer() failed: %v", err)
		return
	}
	defer pd.Close()

	v := struct {
		Data string
	}{}

	v.Data = "Sync Hello, World!"
	err = pd.Produce("test", []any{v, v}, true)
	if err != nil {
		panic(err)
	}

	v.Data = "Async Hello, World!"
	err = pd.Produce("test", []any{v, v}, false)
	if err != nil {
		panic(err)
	}
}

func TestNewConsumer(t *testing.T) {
	con, err := NewConsumer(c, "")
	if err != nil {
		t.Errorf("NewConsumer() failed: %v", err)
		return
	}
	defer con.Close()

	err = con.Consume([]string{"test"}, 0, sarama.OffsetNewest, func(msg *sarama.ConsumerMessage) {
		log.Printf("Consumer: %s\n", string(msg.Value))
	})
	if err != nil {
		panic(err)
	}

	// Wait for a signal to quit
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm
}

func TestNewConsumerGroup(t *testing.T) {
	con, err := NewConsumer(c, "test")
	if err != nil {
		t.Errorf("NewConsumerGroup() failed: %v", err)
		return
	}
	defer con.Close()

	err = con.Consume([]string{"test"}, 0, sarama.OffsetNewest, func(msg *sarama.ConsumerMessage) {
		log.Printf("ConsumerGroup: %s\n", string(msg.Value))
	})

	// Wait for a signal to quit
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm
}
