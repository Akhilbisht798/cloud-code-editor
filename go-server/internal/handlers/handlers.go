package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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

type CreateProjectResponse struct {
	Ip string `json:"ip"`
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
	key := fmt.Sprintf("userId-%d/%s/", user.Id, req.Project)
	log.Println(key)
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
		log.Println("Project Already exits. Starting up the container.")
	} else {
		log.Println("Creating a new Project.")
		err := cloud.CreateNewProjectS3(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	ip, err := cloud.CreateECSContainer(strconv.Itoa(user.Id), req.Project)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	output := CreateProjectResponse{
		Ip: ip,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	err = json.NewEncoder(w).Encode(&output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
		http.Error(w, db.Error.Error(), http.StatusNotFound)
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
		http.Error(w, db.Error.Error(), http.StatusUnauthorized)
		return
	}

	jsonData, err := json.Marshal(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

type S3PresignedUrlRequest struct {
	EmailId   string `json:"userId"`
	ProjectId string `json:"projectId"`
}

func S3PresignedGetURLHandler(w http.ResponseWriter, r *http.Request) {
	var req S3PresignedUrlRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	prefix := fmt.Sprintf("%s/%s", req.EmailId, req.ProjectId)
	var svc cloud.S3Client
	if os.Getenv("APP_ENV") == "production" {
		svc = cloud.GetS3Client()
	} else {
		svc = cloud.GetS3ClientDevelopment()
	}
	cloud.GetPresignerClient()
	urls, err := svc.GetPresignedUrls(os.Getenv("BUCKET"), prefix)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(urls)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)
}
