package repository

import (
	"github.com/jbholae/golang-chat/models"
	"gorm.io/gorm"
)

type User interface {
	GetId() string
	GetName() string
}

type IUserRepository interface {
	AddUser(user User)
	RemoveUser(user User)
	FindUserById(ID string) User
	GetAllUsers() []User
}

type UserRepository struct {
	// Db *sql.DB
	Db *gorm.DB
}

func (repo *UserRepository) AddUser(user models.User) error {
	// stmt, err := repo.Db.Prepare("INSERT INTO user(id, name) values(?,?)")
	// checkErr(err)

	// _, err = stmt.Exec(user.GetId(), user.GetName())
	// checkErr(err)
	return repo.Db.Create(user).Error
}

func (repo *UserRepository) RemoveUser(user models.User) error {
	// stmt, err := repo.Db.Prepare("DELETE FROM user WHERE id = ?")
	// checkErr(err)

	// _, err = stmt.Exec(user.GetId())
	// checkErr(err)
	return repo.Db.Model(&user).Where("id = ?", user.GetId).Error
}

func (repo *UserRepository) FindUserById(ID string) (user models.User) {
	// row := repo.Db.QueryRow("SELECT id, name FROM user where id = ? LIMIT 1", ID)
	row := repo.Db.Raw("SELECT id, name FROM user WHERE id = ? LIMIT 1", ID)

	if err := row.Scan(&user); err != nil {
		// if err == sql.ErrNoRows {
		// 	return nil
		// }
		panic(err)
	}

	return user
}

func (repo *UserRepository) GetAllUsers() (users []models.User) {
	// rows, err := repo.Db.Query("SELECT id, name FROM user")
	rows := repo.Db.Raw("SELECT * FROM user")
	rows.Find(&users)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	var user User
	// 	rows.Scan(&user.Id, &user.Name)
	// 	users = append(users, &user)
	// }

	return users
}
