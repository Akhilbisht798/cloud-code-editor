package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"

	"github.com/gorilla/websocket"
)

var commandChan = make(chan string)

func readCommand(msg Message) {
	command := msg.Data["command"].(string)
	command = fmt.Sprintf("%s\n", command)
	commandChan <- command
}

func runCommand(conn *websocket.Conn) {
	cmd := exec.Command("bash")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println("error creating stdin pipe: ", err)
		return
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("error creating stdout pipe: ", err)
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("error creating stderr pipe: ", err)
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting command: ", err)
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()

			resp := Message{
				Event: "command-response",
				Data: map[string]interface{}{
					"response": line,
				},
			}
			jsonData, err := json.Marshal(resp)
			if err != nil {
				log.Fatal("error converting it to json: ", err)
				continue
			}
			err = conn.WriteMessage(websocket.TextMessage, jsonData)
			if err != nil {
				log.Fatal("Error writing to websocket: ", err)
				break
			}
		}

		if err := scanner.Err(); err != nil {
			log.Println("error reading stdout: ", err)
		}
		stdout.Close()
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()

			resp := Message{
				Event: "command-response",
				Data: map[string]interface{}{
					"response": line,
				},
			}
			jsonData, err := json.Marshal(resp)
			if err != nil {
				log.Fatal("error converting it to json: ", err)
				continue
			}
			err = conn.WriteMessage(websocket.TextMessage, jsonData)
			if err != nil {
				log.Fatal("Error writing to websocket: ", err)
				break
			}
		}

		if err := scanner.Err(); err != nil {
			log.Println("error reading stdout: ", err)
		}

		stderr.Close()
	}()

	go func() {
		for {
			command := <-commandChan
			_, err := stdin.Write([]byte(command))
			if err != nil {
				log.Println("error writing to stdin: ", err)
			}
		}
	}()
}
