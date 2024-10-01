package database

import (
	"github.com/Akhilbisht798/cloud-text-editor/go-server/internal/types"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DbConnect() {
	// sql := "root:rootpassword@tcp(localhost:3306)/db?charset=utf8mb4&parseTime=True&loc=Local"
	dbUrl := ""
	db, err := gorm.Open(
		postgres.Open(dbUrl),
		&gorm.Config{})

	if err != nil {
		panic("error connection to db")
	}
	db.AutoMigrate(&types.User{})
	DB = db
}
