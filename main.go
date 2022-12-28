package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/jbholae/golang-chat/config"
)

var addr = flag.String("addr", ":8000", "http server address")

func main() {
	fmt.Print("on the go")
	flag.Parse()

	// config.CreateRedisClient()

	config.InitDB()
	// defer db.Close()

	// userRepository := &repository.UserRepository{Db: db}

	// wsServer := NewWebsocketServer(repository.RoomRepository{Db: db}, repository.UserRepository{Db: db})
	wsServer := NewWebsocketServer()
	go wsServer.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(wsServer, w, r)
		fmt.Print(" running")
	})

	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	// log.Fatal(http.ListenAndServe(*addr, nil))
	log.Fatal(http.ListenAndServe(":8000", nil))
}
