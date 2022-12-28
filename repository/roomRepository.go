package repository

import (
	"github.com/jbholae/golang-chat/models"
	"gorm.io/gorm"
)

type Room struct {
	Id      string
	Name    string
	Private bool
}

func (room *Room) GetId() string {
	return room.Id
}

func (room *Room) GetName() string {
	return room.Name
}

func (room *Room) GetPrivate() bool {
	return room.Private
}

type RoomRepository struct {
	Db *gorm.DB
}

func (repo *RoomRepository) AddRoom(room models.Room) error {
	// stmt, err := repo.Db.Prepare("INSERT INTO room(id, name, private) values(?,?,?)")
	// checkErr(err)

	// _, err = stmt.Exec(room.GetId(), room.GetName(), room.GetPrivate())
	// checkErr(err)
	return repo.Db.Create(room).Error
}

func (repo *RoomRepository) FindRoomByName(name string) models.Room {

	row := repo.Db.Raw("SELECT id, name, private FROM room where name = ? LIMIT 1", name)
	// row := repo.Db.q

	var room Room

	if err := row.Scan(&room.Name); err != nil {
		// if err == sql.ErrNoRows {
		// 	return nil
		// }
		panic(err)
	}

	return &room

}

// func checkErr(err error) {
// 	if err != nil {
// 		panic(err)
// 	}
// }
