package kafka

import (
	"context"

	"github.com/IBM/sarama"
	cmap "github.com/orcaman/concurrent-map/v2"
)

type Consumer interface {
	Consume(topics []string, partition int32, offset int64, f func(*sarama.ConsumerMessage)) error
	PauseAll()
	ResumeAll()
	Close() error
}

type consumer struct {
	c       sarama.Consumer
	cps     cmap.ConcurrentMap[string, sarama.PartitionConsumer]
	groupID string
	cg      sarama.ConsumerGroup
}

func (c *consumer) isGroup() bool {
	return c.groupID != ""
}

func (c *consumer) Consume(topics []string, partition int32, offset int64, f func(*sarama.ConsumerMessage)) error {
	if len(topics) == 0 {
		return ErrorEmptyTopic
	}

	var err error
	// Consumer
	if !c.isGroup() {
		for i := range topics {
			var cp sarama.PartitionConsumer
			cp, err = c.c.ConsumePartition(topics[i], partition, offset)
			if err != nil {
				return err
			}

			go func() {
				for m := range cp.Messages() {
					f(m)
				}
				for err = range cp.Errors() {
					return
				}
			}()

			c.cps.Set(topics[i], cp)
		}
	}

	// ConsumerGroup
	if c.isGroup() {
		h := new(ConsumerGroupHandler)
		h.f = f
		err = c.cg.Consume(context.TODO(), topics, h)
		go func() {
			for err = range c.cg.Errors() {
				return
			}
		}()
	}

	return err
}

func (c *consumer) PauseAll() {
	if c.groupID == "" {
		c.c.PauseAll()
	}

	if c.groupID != "" {
		c.cg.PauseAll()
	}
}

func (c *consumer) ResumeAll() {
	if c.groupID == "" {
		c.c.ResumeAll()
	}

	if c.groupID != "" {
		c.cg.ResumeAll()
	}
}

func (c *consumer) Close() error {
	var err error
	if c.groupID == "" {
		for _, v := range c.cps.Items() {
			v.AsyncClose()
			err = v.Close()
		}
		err = c.c.Close()
	}

	if c.groupID != "" {
		err = c.cg.Close()
	}

	return err
}

func NewConsumer(c Config, groupID string) (Consumer, error) {
	con := new(consumer)
	con.groupID = groupID
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	switch c.ConsumerGroupStrategy {
	case ConsumerGroupRangeBalanceStrategy:
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	case ConsumerGroupRoundRobinBalanceStrategy:
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	case ConsumerGroupStickyBalanceStrategy:
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	default:
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	}

	var err error
	if !con.isGroup() {
		con.c, err = sarama.NewConsumer(c.Addresses, config)
		con.cps = cmap.New[sarama.PartitionConsumer]()
	}

	if con.isGroup() {
		con.groupID = groupID
		con.cg, err = sarama.NewConsumerGroup(c.Addresses, groupID, config)
	}

	return con, err
}

type ConsumerGroupHandler struct {
	f func(*sarama.ConsumerMessage)
}

func (h *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case m := <-claim.Messages():
			h.f(m)
			sess.MarkMessage(m, "")
		}
	}
}
