package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func dbConnect() {
	sql := "root:rootpassword@tcp(localhost:3306)/db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(
		mysql.Open(sql),
		&gorm.Config{})

	if err != nil {
		panic("error connection to db")
	}
	db.AutoMigrate(&User{})
	DB = db
}
