package models

import "gorm.io/gorm"

type UserRoom struct {
	gorm.Model
	ID     int32   `json:"id" gorm:"primary_key"`
	UserID *User   `json:"user"`
	RoomID []*Room `json:"room" gorm:"foreignkey:Id"`
}
