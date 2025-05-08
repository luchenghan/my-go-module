package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/rs/xid"
)

type Producer interface {
	Produce(topic string, values []any, sync bool) error
	Close() error
}

type producer struct {
	async sarama.AsyncProducer
	sync  sarama.SyncProducer
}

func (p *producer) Close() error {
	err := p.sync.Close()
	err = p.async.Close()
	return err
}

func (p *producer) Produce(topic string, values []any, sync bool) error {
	if topic == "" {
		return ErrorEmptyTopic
	}

	var err error
	msgs := make([]*sarama.ProducerMessage, 0, len(values))
	for _, v := range values {
		var marshal []byte
		marshal, err = json.Marshal(v)
		if err != nil {
			return err
		}

		msg := new(sarama.ProducerMessage)
		msg.Topic = topic
		msg.Value = sarama.ByteEncoder(marshal)
		msgs = append(msgs, msg)
	}

	if sync && !p.sync.IsTransactional() {
		if err = p.sync.SendMessages(msgs); err != nil {
			return err
		}
	}

	if sync && p.sync.IsTransactional() {
		if err = p.sync.BeginTxn(); err != nil {
			return err
		}

		if err = p.sync.SendMessages(msgs); err != nil {
			return err
		}

		if err = p.sync.CommitTxn(); err != nil {
			return err
		}
	}

	if !sync && p.async.IsTransactional() {
		if err = p.async.BeginTxn(); err != nil {
			return err
		}

		for _, msg := range msgs {
			p.async.Input() <- msg
		}

		for i := 0; i < len(msgs); {
			select {
			case _ = <-p.async.Successes():
				i++
			case err = <-p.async.Errors():
				return err
			}
		}

		if err = p.async.CommitTxn(); err != nil {
			return err
		}
	}

	if !sync && !p.async.IsTransactional() {
		for _, msg := range msgs {
			p.async.Input() <- msg
		}
	}

	return err
}

type Option func(config *sarama.Config)

func WithTxn() Option {
	return func(c *sarama.Config) {
		c.Producer.Transaction.ID = fmt.Sprintf("tx-%v", xid.New().String())
		c.Producer.Idempotent = true
		c.Producer.RequiredAcks = sarama.WaitForAll
		c.Consumer.IsolationLevel = sarama.ReadCommitted
		c.Net.MaxOpenRequests = 1
	}
}

func WithRequiredAcks(r sarama.RequiredAcks) Option {
	return func(c *sarama.Config) {
		c.Producer.RequiredAcks = r
	}
}

func NewProducer(c Config, options ...Option) (Producer, error) {
	p := new(producer)
	syncConfig := sarama.NewConfig()
	asyncConfig := sarama.NewConfig()
	syncConfig.Producer.Return.Successes = true
	asyncConfig.Producer.Return.Successes = true

	switch c.ProducerPartitioner {
	case ProducerPartitionHash:
		syncConfig.Producer.Partitioner = sarama.NewHashPartitioner
		asyncConfig.Producer.Partitioner = sarama.NewHashPartitioner
	case ProducerPartitionRef:
		syncConfig.Producer.Partitioner = sarama.NewReferenceHashPartitioner
		asyncConfig.Producer.Partitioner = sarama.NewReferenceHashPartitioner
	case ProducerPartitionCons:
		syncConfig.Producer.Partitioner = sarama.NewConsistentCRCHashPartitioner
		asyncConfig.Producer.Partitioner = sarama.NewConsistentCRCHashPartitioner
	default:
		syncConfig.Producer.Partitioner = sarama.NewHashPartitioner
		asyncConfig.Producer.Partitioner = sarama.NewHashPartitioner
	}

	for _, option := range options {
		option(asyncConfig)
		option(syncConfig)
	}

	if err1, err2 := syncConfig.Validate(), asyncConfig.Validate(); err1 != nil || err2 != nil {
		return nil, fmt.Errorf("syncConfig.Validate() failed: %v, asyncConfig.Validate() failed: %v", err1, err2)
	}

	var err error
	if p.async, err = sarama.NewAsyncProducer(c.Addresses, asyncConfig); err != nil {
		return nil, err
	}

	if p.sync, err = sarama.NewSyncProducer(c.Addresses, syncConfig); err != nil {
		return nil, err
	}

	return p, nil
}
