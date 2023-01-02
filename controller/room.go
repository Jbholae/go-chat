package controller

import (
	"encoding/json"
	"fmt"
	"golang-chat/models"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type RoomController struct {
	db *gorm.DB
}

func NewRoomController(db *gorm.DB) RoomController {
	return RoomController{
		db: db,
	}
}

func (c RoomController) GetRooms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var userRoom []models.Room
	params := r.URL.Query()
	userId, _ := strconv.ParseInt(params["UserId"][0], 10, 64)

	query := c.db.Model(&models.Room{}).
		Joins("JOIN user_rooms on rooms.id = user_rooms.room_id").
		Where("user_rooms.user_id = ?", userId).
		Find(&userRoom).
		Error

	if query != nil {
		log.Println("Error finding rooms", query.Error())
		json.NewEncoder(w).Encode(map[string]string{
			"message": fmt.Sprintf("error finding rooms of user :: %v", userId),
		})
		return
	}

	json.NewEncoder(w).Encode(userRoom)
}
