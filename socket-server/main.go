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
	godotenv.Load()
	SERVER = os.Getenv("SERVER_URL")
	userId = os.Getenv("userId")
	projectId = os.Getenv("projectId")
	if SERVER == "" {
		fmt.Println("server url not set")
		return
	}
	SERVER = fmt.Sprintf("%s:%s", SERVER, "8080")

	filePath := fmt.Sprintf("%s/%s", userId, projectId)
	fmt.Println(filePath)
	err := getFilesFromS3(userId, projectId)
	if err != nil {
		fmt.Println("error", err.Error())
		err = saveDir(filePath)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	log.Print("Server Starting")
	log.SetFlags(0)
	http.HandleFunc("/ws", echo)
	log.Fatal(http.ListenAndServe(":5000", nil))
}
