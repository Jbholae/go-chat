package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jbholae/golang-chat/models"
	"gorm.io/gorm"
)

var db *gorm.DB

func CreatedUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "aplication/json")
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)
	db.Create(&user)
	json.NewEncoder(w).Encode(user)
}

func GetRooms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var userRoom []models.UserRoom
	params := mux.Vars(r)
	db.Preload("RoomID").Find(&userRoom, params["UserID"])
	json.NewEncoder(w).Encode(userRoom)

}
