package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

type Message struct {
	Event string                 `json:"event"`
	Data  map[string]interface{} `json:"data"`
}

var upgrader = websocket.Upgrader{
	//Allow all origin
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	runCommand(c)
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		if mt == websocket.TextMessage {
			eventHandler(c, message)
		}
	}
}

func eventHandler(conn *websocket.Conn, message []byte) {
	var msg Message
	err := json.Unmarshal(message, &msg)
	if err != nil {
		fmt.Println("Error decoding JSON: ", err)
		return
	}

	switch msg.Event {
	case "command":
		readCommand(msg)
	case "sendFile":
		sendFilesToClient(conn, msg)
	case "Filechanges":
		fileChanges(msg)
	}
}

// Handler Functions.
func fileChanges(msg Message) {
	path, ok := msg.Data["file"].(string)
	if !ok {
		log.Println("invalid", ok)
		return
	}
	content, ok := msg.Data["content"].(string)
	if !ok {
		log.Println("invalid", ok)
		return
	}
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		log.Fatal("Error writing to file: ", err)
	}
}

func sendFilesToClient(conn *websocket.Conn, msg Message) {
	path, ok := msg.Data["data"].(string)
	if !ok {
		log.Println("invalid path: ", ok)
		return
	}
	files, err := readDir(path)
	if err != nil {
		return
	}
	f, err := mapToJson(files)
	if err != nil {
		log.Println(err)
		return
	}
	conn.WriteMessage(websocket.TextMessage, []byte(f))
}
