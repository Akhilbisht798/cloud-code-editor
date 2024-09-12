package database

import (
	"github.com/Akhilbisht798/cloud-text-editor/go-server/internal/types"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DbConnect() {
	sql := "root:rootpassword@tcp(localhost:3306)/db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(
		mysql.Open(sql),
		&gorm.Config{})

	if err != nil {
		panic("error connection to db")
	}
	db.AutoMigrate(&types.User{})
	DB = db
}
