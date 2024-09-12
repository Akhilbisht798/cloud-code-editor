package handlers

import (
	"net/http"
	"os"

	"github.com/Akhilbisht798/cloud-text-editor/go-server/internal/database"
	"github.com/Akhilbisht798/cloud-text-editor/go-server/internal/types"
	"github.com/dgrijalva/jwt-go/v4"
)

var SECRET = os.Getenv("SECRET_KEY")

func authHandler(r *http.Request) (types.User, error) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return types.User{}, err
	}
	token, err := jwt.ParseWithClaims(cookie.Value, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET), nil
	})

	if err != nil {
		return types.User{}, nil
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user types.User
	db := database.DB.Where("id = ?", claims.Issuer).First(&user)
	if db.Error != nil {
		return types.User{}, nil
	}
	return user, nil
}
