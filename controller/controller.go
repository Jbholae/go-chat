package controller

import (
	"encoding/json"
	"go.uber.org/fx"
	"log"
	"net/http"
)

var Module = fx.Options(
	fx.Provide(NewUserController),
	fx.Provide(NewRoomController),
)

// GetNothing :: Prints Doing nothing
func GetNothing(w http.ResponseWriter, r *http.Request) {
	log.Println("Doing nothing")
	json.NewEncoder(w).Encode("asdaasdas")
}
