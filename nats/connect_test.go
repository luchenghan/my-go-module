package nats

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
)

func TestNewNatsConnection(t *testing.T) {
	natsUrl := "nats://127.0.0.1:4222,nats://127.0.0.1:4223,nats://127.0.0.1:4224"
	nc, err := NewNatsConnection(natsUrl)
	if err != nil {
		log.Fatal(err)
	}

	nc.Close()
}

func TestNatsReplyRequest(t *testing.T) {
	natsUrl := "nats://127.0.0.1:4222,nats://127.0.0.1:4223,nats://127.0.0.1:4224"
	nc, err := NewNatsConnection(natsUrl)
	if err != nil {
		log.Fatal(err)
	}

	// Replies
	nc.Subscribe("help", func(m *nats.Msg) {
		nc.Publish(m.Reply, []byte("I can help!"))
	})
	fmt.Println(nc.NumSubscriptions())

	// Requests
	msg, err := nc.Request("help", []byte("help me"), 10*time.Millisecond)
	fmt.Println(string(msg.Data))

	nc.Close()
}
