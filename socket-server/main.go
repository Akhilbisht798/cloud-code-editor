package main

import (
	//"flag"
	"fmt"
	"log"
	"net/http"
)

// var addr = flag.String("addr", "localhost:8080", "http service address")
const SERVER = "http://localhost:3000"

var userId string
var projectId string

func main() {
	// userId = os.Getenv("userId")
	// projectId = os.Getenv("projectId")
	userId = "1"
	projectId = "client"

	filePath := fmt.Sprintf("%s/%s", userId, projectId)
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
	//flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/ws", echo)
	//log.Fatal(http.ListenAndServe(*addr, nil))
	log.Fatal(http.ListenAndServe(":5000", nil))
}
