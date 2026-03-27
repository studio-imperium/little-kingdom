package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	go create_client(conn)
}

func main() {
	fmt.Println("Game running on 8082")
	http.HandleFunc("/game", handler)
	log.Fatal(http.ListenAndServe(":8082", nil))
}
