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
		stdin.Close()
		stdout.Close()
		stderr.Close()
		err := cmd.Process.Kill()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}()

	writeToSocket := func(scanner *bufio.Scanner, pipeName string) {
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
				log.Println("error converting it to json: ", err)
				continue
			}

			err = conn.WriteMessage(websocket.TextMessage, jsonData)
			if err != nil {
				log.Println("Error writing to websocket: ", err)
				break
			}
		}
		if err := scanner.Err(); err != nil {
			log.Printf("error reading %s: %v", pipeName, err)
		}
	}

	go func() {
		defer stdout.Close()
		scanner := bufio.NewScanner(stdout)
		writeToSocket(scanner, "stdout")
	}()

	go func() {
		defer stderr.Close()
		scanner := bufio.NewScanner(stderr)
		writeToSocket(scanner, "stderr")
	}()

	go func() {
		defer stdin.Close()
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
