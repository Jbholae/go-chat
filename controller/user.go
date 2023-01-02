package controller

import (
	"encoding/json"
	"golang-chat/models"
	"gorm.io/gorm"
	"net/http"
)

type UserController struct {
	db *gorm.DB
}

func NewUserController(
	db *gorm.DB,
) UserController {
	return UserController{
		db: db,
	}
}

func (c UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "aplication/json")
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	c.db.Create(&user)

	json.NewEncoder(w).Encode(user)
}
