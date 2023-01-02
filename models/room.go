package models

import (
	"gorm.io/gorm"
)

type Room struct {
	gorm.Model
	Name    string
	Private bool
}

func (room *Room) GetId() uint {
	return room.Model.ID
}

func (room *Room) GetName() string {
	return room.Name
}

func (room *Room) GetPrivate() bool {
	return room.Private
}
