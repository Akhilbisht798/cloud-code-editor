package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var SERVER string
var userId string
var projectId string

func main() {
	log.Println("Starting Server")
	godotenv.Load()
	SERVER = os.Getenv("SERVER_URL")
	userId = os.Getenv("userId")
	projectId = os.Getenv("projectId")
	if SERVER == "" {
		fmt.Println("server url not set")
		return
	}
	SERVER = fmt.Sprintf("http://%s:%s", SERVER, "8080")

	filePath := fmt.Sprintf("%s/%s", userId, projectId)
	log.Println("filepath: ", filePath)
	err := getFilesFromS3(userId, projectId)
	if err != nil {
		fmt.Println("error", err.Error())
		return
		// err = saveDir(filePath)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
	}
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/ws", echo)
	log.Println("Server Listing at port: ", ":5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}
