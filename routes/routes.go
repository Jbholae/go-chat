package routes

import (
	"github.com/gorilla/mux"
	"github.com/jbholae/golang-chat/controller"
)

func initializeRouter() {
	r := mux.NewRouter()

	r.HandleFunc("api/user", controller.CreatedUser).Methods("POST")
	r.HandleFunc("api/room", controller.GetRooms).Methods("Get")

	// log.Fatal(http.ListenAndServe(":9000", handlers.CORSA.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r))
}
