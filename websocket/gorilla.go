package websocket

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebSocketHandleFunc(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("upgrader.Upgrade err: %v\n", err)
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("conn.ReadMessage err: %v\n", err)
			break
		}
		fmt.Printf("Received message: %s\n", message)

		// 回显消息
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			fmt.Printf("conn.WriteMessage err: %v\n", err)
			break
		}
	}
}

var hub = NewHub()

func init() {
	go hub.Run()
}

func WebSocketHandleFuncWithHub(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("upgrader.Upgrade err: %v\n", err)
		return
	}
	defer conn.Close()

	client := &Client{conn: conn, send: make(chan []byte)}

	hub.register <- client

	defer func() {
		hub.unregister <- client
		conn.Close()
	}()

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("client.ReadMessage err: %v\n", err)
				hub.unregister <- client
				break
			}
			log.Printf("Received message: %s\n", message)
			client.send <- message
		}
	}()

	for message := range client.send {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("client.WriteMessage err: %v\n", err)
			return
		}
	}
}
