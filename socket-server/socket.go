package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"

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
	case "sendData":
		sendData(conn, msg)
	case "command":
		executeCommand(conn, msg)
	case "sendFile":
		sendFilesToClient(conn, msg)
	}
}

// Handler Functions.
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
	f, err := mapToString(files)
	if err != nil {
		log.Println(err)
		return
	}
	conn.WriteMessage(websocket.TextMessage, []byte(f))
}

func executeCommand(conn *websocket.Conn, msg Message) {
	command, ok := msg.Data["data"].(string)
	if !ok {
		log.Println("Invalid command")
		return
	}

	command = strings.TrimSpace(command)
	if command == "" {
		errMsg := "empty command"
		conn.WriteMessage(websocket.TextMessage, []byte(errMsg))
		return
	}

	cmdParts := strings.Fields(command)
	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error Executing command: %s\n", err)
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error: %sOutput: %s", err, output)))
		return
	}

	conn.WriteMessage(websocket.TextMessage, output)
}

func sendData(conn *websocket.Conn, msg Message) {
	name, ok := msg.Data["name"].(string)
	if !ok {
		log.Println("Invalid name data")
		return
	}
	err := conn.WriteMessage(websocket.TextMessage, []byte(name))
	if err != nil {
		log.Println("error sending message: ", err)
		return
	}
}
