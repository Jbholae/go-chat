package config

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// server_port := 8000;
var db_host = "localhost"
var db_port = "3306"
var db_name = "chat"
var db_username = "root"
var db_password = "Mysql@123"

func InitDB() *gorm.DB {
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", db_username, db_password, db_host, db_port, db_name)

	db, err := gorm.Open(mysql.Open(url), &gorm.Config{})
	_ = db.Exec("CREATE DATABASE IF NOT EXISTS " + db_name + ";")

	// db, err := sql.Open("sqlite3", "./chatdb.db")
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate()

	// sqlStmt := `
	// CREATE TABLE IF NOT EXISTS room (
	// 	id VARCHAR(255) NOT NULL PRIMARY KEY,
	// 	name VARCHAR(255) NOT NULL,
	// 	private TINYINT NULL
	// );
	// `
	// _, err = db.Exec(sqlStmt)
	// if err != nil {
	// 	log.Fatal("%q: %s\n", err, sqlStmt)
	// }

	// sqlStmt = `
	// CREATE TABLE IF NOT EXISTS user (
	// 	id VARCHAR(255) NOT NULL PRIMARY KEY,
	// 	name VARCHAR(255) NOT NULL
	// );
	// `
	// _, err = db.Exec(sqlStmt)
	// if err != nil {
	// 	log.Fatal("%q: %s\n", err, sqlStmt)
	// }

	return db
}
