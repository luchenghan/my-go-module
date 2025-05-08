package websocket

import (
	"log"
	"net/http"
	"testing"

	"github.com/gorilla/websocket"
)

func TestWebSocketClient(t *testing.T) {
	conn1, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		log.Fatal("dial:", err)
		return
	}
	defer conn1.Close()

	//conn2, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	//if err != nil {
	//	log.Fatal("dial:", err)
	//	return
	//}
	//defer conn2.Close()

	err = conn1.WriteMessage(websocket.TextMessage, []byte("hello websocket"))
	if err != nil {
		log.Printf("write: %v", err)
		return
	}

	_, message, err := conn1.ReadMessage()
	if err != nil {
		log.Printf("read: %v", err)
		return
	}

	err = conn1.WriteMessage(websocket.TextMessage, []byte("hi websocket"))
	if err != nil {
		log.Printf("write: %v", err)
		return
	}

	_, message2, err := conn1.ReadMessage()
	if err != nil {
		log.Printf("read: %v", err)
		return
	}

	log.Printf("conn1 recv: %s, conn1 recv: %s", message, message2)

	//
	//_, message3, err := conn1.ReadMessage()
	//if err != nil {
	//	log.Printf("read: %v", err)
	//	return
	//}
	//
	//_, message4, err := conn2.ReadMessage()
	//if err != nil {
	//	log.Printf("read: %v", err)
	//	return
	//}
	//
	//log.Printf("conn1 recv: %s, conn2 recv: %s", message3, message4)

}

func TestWebSocketHandleFunc(t *testing.T) {
	http.HandleFunc("/ws", WebSocketHandleFunc)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func TestWebSocketHandleFuncWithHub(t *testing.T) {
	http.HandleFunc("/ws", WebSocketHandleFuncWithHub)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
