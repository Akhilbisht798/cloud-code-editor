package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Akhilbisht798/cloud-text-editor/go-server/internal/cloud"
	"github.com/Akhilbisht798/cloud-text-editor/go-server/internal/database"
	"github.com/Akhilbisht798/cloud-text-editor/go-server/internal/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/dgrijalva/jwt-go/v4"
	"golang.org/x/crypto/bcrypt"
)

type CreateProjectRequest struct {
	Project string `json:"project"`
}

func CreateProject(w http.ResponseWriter, r *http.Request) {
	user, err := authHandler(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	var req CreateProjectRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var svc cloud.S3Client
	if os.Getenv("APP_ENV") == "production" {
		svc = cloud.GetS3Client()
	} else {
		svc = cloud.GetS3ClientDevelopment()
	}
	key := fmt.Sprintf("%s/%s/", user.Email, req.Project)
	fmt.Println(key)
	maxKey := int32(1)

	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(os.Getenv("BUCKET")),
		Prefix:  &key,
		MaxKeys: &maxKey,
	}

	result, err := svc.Client.ListObjectsV2(context.TODO(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(result.Contents) > 0 {
		fmt.Println("ALready exists")
	} else {
		fmt.Println("create a new directory and start with that with readme.md")
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Created a new project"))
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	var data map[string]string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 12)
	user := types.User{
		Email:    data["email"],
		Password: string(password),
	}
	db := database.DB.Create(&user)
	if db.Error != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	jsonData, _ := json.Marshal(user)
	w.Write(jsonData)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	var data map[string]string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	var user types.User

	db := database.DB.Where("email = ?", data["email"]).First(&user)
	if db.Error != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if user.Id == 0 {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"]))
	if err != nil {
		http.Error(w, "wrong password", http.StatusBadRequest)
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer: strconv.Itoa(user.Id),
		ExpiresAt: &jwt.Time{
			Time: time.Now().Add(time.Hour * 24),
		},
	})

	token, err := claims.SignedString([]byte(SECRET))
	if err != nil {
		http.Error(w, "Could'nt login", http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	cookie, err := r.Cookie("jwt")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	token, err := jwt.ParseWithClaims(cookie.Value, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET), nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user types.User
	db := database.DB.Where("id = ?", claims.Issuer).First(&user)
	if db.Error != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	jsonData, err := json.Marshal(&user)
	w.Write([]byte(jsonData))
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    "/",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}
