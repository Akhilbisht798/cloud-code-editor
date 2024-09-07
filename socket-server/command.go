package main

import (
	"bufio"
	"context"
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

func runCommand(ctx context.Context, conn *websocket.Conn) {
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

	//TODO: safely close these go routine and start new ones on reconnection
	go func() {
		<-ctx.Done()
		err := cmd.Process.Kill()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}()

	go func() {
		defer stdout.Close()
		scanner := bufio.NewScanner(stdout)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				scanner.Scan()
				if err := scanner.Err(); err != nil {
					log.Println("error reading stdout: ", err)
				}
				line := scanner.Text()

				resp := Message{
					Event: "command-response",
					Data: map[string]interface{}{
						"response": line,
					},
				}
				jsonData, err := json.Marshal(resp)
				if err != nil {
					log.Println("error converting it to json: ", err)
					continue
				}
				err = conn.WriteMessage(websocket.TextMessage, jsonData)
				if err != nil {
					log.Println("Error writing to websocket: ", err)
					return
				}
			}
		}

	}()

	go func() {
		defer stderr.Close()
		scanner := bufio.NewScanner(stderr)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				scanner.Scan()
				line := scanner.Text()
				if err := scanner.Err(); err != nil {
					log.Println("error reading stdout: ", err)
					return
				}

				resp := Message{
					Event: "command-response",
					Data: map[string]interface{}{
						"response": line,
					},
				}
				jsonData, err := json.Marshal(resp)
				if err != nil {
					log.Println("error converting it to json: ", err)
					continue
				}
				err = conn.WriteMessage(websocket.TextMessage, jsonData)
				if err != nil {
					log.Println("Error writing to websocket: ", err)
					return
				}
			}
		}

	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				command := <-commandChan
				_, err := stdin.Write([]byte(command))
				if err != nil {
					log.Println("error writing to stdin: ", err)
				}
			}
		}
	}()
}
