package server

import (
	"log"

	"github.com/gorilla/websocket"
)

type Packet struct {
	Header string `json:"header"`
}

type Character struct {
	hand      int
	inventory struct {
	}
	Gear struct {
		Head *string
		Body *string
	} `json:"gear"`
}

type Client struct {
	conn           *websocket.Conn
	character_data Character
}

func create_client(conn *websocket.Conn) {

}

func recieve_packets(conn *websocket.Conn) {
	defer conn.Close()
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = conn.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func send_packets(conn *websocket.Conn) {

}
