package types

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"-"`
}
