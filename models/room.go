package models

import (
	"time"

	"gorm.io/gorm"
)

type Room struct {
	gorm.Model
	Id        int64
	Name      string
	Private   bool
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

func (room *Room) GetId() int64 {
	return room.Id
}

func (room *Room) GetName() string {
	return room.Name
}

func (room *Room) GetPrivate() bool {
	return room.Private
}
