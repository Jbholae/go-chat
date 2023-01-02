package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (user *User) GetId() uint     { return user.Model.ID }
func (user *User) GetName() string { return user.Name }
