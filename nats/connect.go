package nats

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

func NewNatsConnection(url string) (*nats.Conn, error) {
	return nats.Connect(url,
		nats.UserInfo("rd", "password"),
		nats.ConnectHandler(func(conn *nats.Conn) {
			fmt.Println(conn.Status().String())
		}),
		nats.DisconnectErrHandler(func(conn *nats.Conn, err error) {
			fmt.Printf("Got disconnected! Reason: %v\n", err)
		}),
		nats.ReconnectHandler(func(conn *nats.Conn) {
			fmt.Printf("Got reconnected to %v!\n", conn.ConnectedUrl())
		}),
		nats.MaxReconnects(5),
		nats.ReconnectWait(2*time.Second),
		nats.ClosedHandler(func(conn *nats.Conn) {
			fmt.Printf("Connection closed. Reason: %v\n", conn.LastError())
		}),
	)
}
