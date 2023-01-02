package models

import "gorm.io/gorm"

type UserRoom struct {
	gorm.Model
	UserID uint `json:"user_id"`
	RoomID uint `json:"room_id"`
}
