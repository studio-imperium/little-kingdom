package main

import (
	"fmt"
	"log"
	"net/http"
	"server/engine"
	"server/packets"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func connector(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	go packets.CreateClient(conn)
}

func main() {
	engine.InitAssets()
	go engine.Worlds[0].Run()
	go packets.PropogateWorldState()

	fmt.Println("Listening on 8082")
	http.HandleFunc("/connect", connector)
	log.Fatal(http.ListenAndServe(":8082", nil))
}
