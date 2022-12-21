package main

import (
	"fmt"
	"net/http"

	"github.com/jbholae/golang-chat/pkg/websocket"
)

func serverWS(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("websocket endpoint reached")

	conn, err := websocket.Upgrade(w, r)

	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &websocket.Client{
		Conn: conn,
		Pool: pool,
	}
	pool.Register <- client
	client.Read()
}

func setupRoutes() {
	pool := websocket.NewPool()
	go pool.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serverWS(pool, w, r)

	})
}

func main() {
	fmt.Println("Chat")
	setupRoutes()
	http.ListenAndServe(":9000", nil)
}
