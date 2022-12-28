package repository

import (
	"github.com/jbholae/golang-chat/models"
	"gorm.io/gorm"
)

type IRoomRepository interface {
	AddRoom(room models.Room)
	FindRoomByName(name string) models.Room
}

type RoomRepository struct {
	Db *gorm.DB
}

func (repo *RoomRepository) AddRoom(room models.Room) error {
	return repo.Db.Create(room).Error
}

func (repo *RoomRepository) FindRoomByName(name string) (room models.Room) {

	row := repo.Db.Raw("SELECT id, name, private FROM room WHERE name = ? LIMIT 1", name)
	// row := repo.Db.q

	if err := row.Scan(&room.Name); err != nil {
		// if err == sql.ErrNoRows {
		// 	return nil
		// }
		panic(err)
	}

	return room
}
