package nats

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/nats-io/nats.go"
)

var c = Config{
	Address: "nats://rd:password@localhost:4222,nats://rd:password@localhost:4223,nats://rd:password@localhost:4224",
}

func TestNatsSyncSubscribe(t *testing.T) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	conn, err := NewConn(c)
	if err != nil {
		panic(err)
	}

	ch := make(chan *nats.Msg)
	err = conn.SyncSubscribe("test", "test", ch)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case msg := <-ch:
			log.Printf("Received message: %s", string(msg.Data))
		case <-sigChan:
			log.Println("Received signal, exiting...")
			if err = conn.Close(); err != nil {
				panic(err)
			}
			return
		}
	}
}

func TestNatsAsyncSubscribeWithCallBack(t *testing.T) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	conn, err := NewConn(c)
	if err != nil {
		panic(err)
	}

	err = conn.AsyncSubscribeWithCallback("test", "", func(msg *nats.Msg) {
		log.Printf("Received message: %s", string(msg.Data))
	})

	select {
	case <-sigChan:
		log.Println("Received signal, exiting...")
	}

	if err = conn.Close(); err != nil {
		panic(err)
	}
}

func TestNatsAsyncSubscribeWithChannel(t *testing.T) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	conn, err := NewConn(c)
	if err != nil {
		panic(err)
	}

	ch := make(chan *nats.Msg)
	err = conn.AsyncSubscribeWithChannel("test", "", ch)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case msg := <-ch:
			log.Printf("Received message: %s", string(msg.Data))
		case <-sigChan:
			log.Println("Received signal, exiting...")
			if err = conn.Close(); err != nil {
				panic(err)
			}
			return
		}
	}
}

func TestDuplicateSubscribe(t *testing.T) {
	conn, err := NewConn(c)
	if err != nil {
		panic(err)
	}

	err = conn.SyncSubscribe("test", "", nil)
	if err != nil {
		panic(err)
	}

	err = conn.SyncSubscribe("test", "", nil)
	if !errors.Is(err, ErrDuplicateSubscription) {
		panic(err)
	}

	err = conn.Close()
	if err != nil {
		panic(err)
	}
}

func TestNatsPublish(t *testing.T) {
	conn, err := NewConn(c)
	if err != nil {
		panic(err)
	}

	err = conn.Publish("test", "", []byte("hello"))
	if err != nil {
		panic(err)
	}

	err = conn.Close()
	if err != nil {
		panic(err)
	}
}
