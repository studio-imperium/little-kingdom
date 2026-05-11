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

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	engine.InitAssets()
	go engine.Worlds[0].Run()
	go packets.PropogateWorldState()

	fmt.Println("Listening on 8082")
	http.HandleFunc("/connect", connector)
	http.Handle(
		"/assets/",
		withCORS(http.StripPrefix("/assets/", http.FileServer(http.FS(engine.JSONAssets())))),
	)
	log.Fatal(http.ListenAndServe(":8082", nil))
}
