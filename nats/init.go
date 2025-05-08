package nats

import (
	"errors"
	"time"

	"github.com/nats-io/nats.go"
)

var ErrDuplicateSubscription = errors.New("duplicate subscription")

const (
	syncSubscribeNextMsgTimeout = 100 * time.Millisecond
)

type Nats interface {
	Publish(subj, reply string, data []byte) error
	SyncSubscribe(subj, queue string, ch chan *nats.Msg) error
	AsyncSubscribeWithCallback(subj, queue string, cb nats.MsgHandler) error
	AsyncSubscribeWithChannel(subj, queue string, ch chan *nats.Msg) error
	Unsubscribe(subj string) error
	Close() error
}

type Config struct {
	Address string `json:"address"` // e.g. "nats://user:password@localhost:4222,nats://user:password@localhost:4223,nats://user:password@localhost:4224" "
}

type nc struct {
	*nats.Conn
	subscriptions map[string]*nats.Subscription
}

func (n *nc) Unsubscribe(subj string) error {
	err := n.subscriptions[subj].Drain()
	if err != nil {
		return err
	}

	delete(n.subscriptions, subj)

	return nil
}

func (n *nc) Close() error {
	for _, sub := range n.subscriptions {
		if err := n.Unsubscribe(sub.Subject); err != nil {
			return err
		}
	}

	if err := n.Conn.Drain(); err != nil {
		return err
	}

	n.Conn.Close()

	return nil
}

// SyncSubscribe subscribes to a subject with a queue group and a channel.
// If queue is empty, it will subscribe to the subject without a queue group.
// If queue is not empty, it will subscribe to the subject with a queue group.
func (n *nc) SyncSubscribe(subj, queue string, ch chan *nats.Msg) error {
	var sub *nats.Subscription
	var err error

	if _, ok := n.subscriptions[subj]; ok {
		return ErrDuplicateSubscription
	}

	if queue != "" {
		sub, err = n.Conn.QueueSubscribeSync(subj, queue)
		n.subscriptions[subj] = sub
	}

	if queue == "" {
		sub, err = n.Conn.SubscribeSync(subj)
		n.subscriptions[subj] = sub
	}

	if err != nil {
		return err
	}

	go func() error {
		for {
			var msg *nats.Msg
			var pMsgs, pBytes int
			pMsgs, pBytes, err = sub.Pending()
			if err != nil {
				return err
			}

			if pMsgs == 0 || pBytes == 0 {
				time.Sleep(syncSubscribeNextMsgTimeout)
				continue
			}

			msg, err = sub.NextMsg(syncSubscribeNextMsgTimeout)
			if err != nil {
				return err
			}

			ch <- msg
		}
	}()
	if err != nil {
		return err
	}

	return nil
}

// Publish publishes a message to a subject.
// If reply is not empty, it will publish a request to a subject.
func (n *nc) Publish(subj, reply string, data []byte) error {
	if reply != "" {
		return n.Conn.PublishRequest(subj, reply, data)
	}

	return n.Conn.Publish(subj, data)
}

// SubscribeWithChannel async subscribes to a subject with a queue group and a channel.
// If queue is empty, it will subscribe to the subject without a queue group.
// If queue is not empty, it will subscribe to the subject with a queue group.
func (n *nc) AsyncSubscribeWithChannel(subj, queue string, ch chan *nats.Msg) error {
	var sub *nats.Subscription
	var err error

	if _, ok := n.subscriptions[subj]; ok {
		return ErrDuplicateSubscription
	}

	if queue != "" {
		sub, err = n.Conn.QueueSubscribeSyncWithChan(subj, queue, ch)
		if err != nil {
			return err
		}
		n.subscriptions[subj] = sub
	}

	if queue == "" {
		sub, err = n.Conn.ChanSubscribe(subj, ch)
		if err != nil {
			return err
		}
		n.subscriptions[subj] = sub
	}

	return nil
}

// AsyncSubscribeWithCallBack subscribes to a subject with a queue group and a callback.
// If queue is empty, it will subscribe to the subject without a queue group.
// If queue is not empty, it will subscribe to the subject with a queue group.
func (n *nc) AsyncSubscribeWithCallback(subj, queue string, cb nats.MsgHandler) error {
	var sub *nats.Subscription
	var err error

	if _, ok := n.subscriptions[subj]; ok {
		return ErrDuplicateSubscription
	}

	if queue != "" {
		sub, err = n.Conn.QueueSubscribe(subj, queue, cb)
		if err != nil {
			return err
		}

		n.subscriptions[subj] = sub
	}

	if queue == "" {
		sub, err = n.Conn.Subscribe(subj, cb)
		if err != nil {
			return err
		}

		n.subscriptions[subj] = sub
	}

	return nil
}

func NewConn(c Config) (Nats, error) {
	n := new(nc)
	var err error

	n.Conn, err = nats.Connect(c.Address)
	if err != nil {
		return nil, err
	}

	n.subscriptions = make(map[string]*nats.Subscription)

	return n, nil
}
