package main

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"-"`
}

type Project struct {
	Id     int `json:"id"`
	UserId int `json:"userId"`
}
