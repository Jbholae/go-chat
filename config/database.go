package config

import (
	"fmt"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"

	"golang-chat/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// server_port := 8000;
var db_host = "localhost"
var db_port = "8889"
var db_name = "chat"
var db_username = "root"
var db_password = "root"

func InitDB() *gorm.DB {
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", db_username, db_password, db_host, db_port, db_name)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)
	db, err := gorm.Open(mysql.Open(url), &gorm.Config{Logger: newLogger})
	_ = db.Exec("CREATE DATABASE IF NOT EXISTS " + db_name + ";")

	// db, err := sql.Open("sqlite3", "./chatdb.db")
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&models.User{}, &models.Room{}, &models.UserRoom{})

	return db
}
