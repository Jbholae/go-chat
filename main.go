package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jbholae/golang-chat/controller"
)

func main() {
	// config.CreateRedisClient()

	//_ = config.InitDB()
	// defer db.Close()

	// userRepository := &repository.UserRepository{Db: db}

	// wsServer := NewWebsocketServer(repository.RoomRepository{Db: db}, repository.UserRepository{Db: db})
	r := mux.NewRouter()

	r.HandleFunc("api/user", controller.CreatedUser).Methods("POST")
	r.HandleFunc("api/room", controller.GetRooms).Methods("Get")
	
	wsServer := NewWebsocketServer()
	go wsServer.Run()

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(wsServer, w, r)
		fmt.Println("running")
	})

	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	// log.Fatal(http.ListenAndServe(*addr, nil))
	log.Fatal(http.ListenAndServe(":8000", nil))
}
