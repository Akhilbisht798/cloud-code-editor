package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

type FileInfo struct {
	Path    string `json:"path"`
	IsDir   bool   `json:"isDir"`
	Content string `json:"content,omitempty"`
}

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
	path, ok := msg.Data["path"].(string)
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

func readDir(dirpath string) ([]FileInfo, error) {
	files := make([]FileInfo, 0)
	entries, err := os.ReadDir(dirpath)
	if err != nil {
		log.Fatal("error reading directory: ", err)
		return nil, err
	}
	for _, e := range entries {
		filePath := dirpath + "/" + e.Name()
		path := dirpath + e.Name()
		f := FileInfo{Path: path, IsDir: e.IsDir()}
		if !e.IsDir() {
			fileContent, err := readFile(filePath)
			if err != nil {
				return nil, err
			}
			f.Content = string(fileContent)
		}
		files = append(files, f)
	}

	return files, nil
}

func readFile(filePath string) ([]byte, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal("Error reading file: ", err)
		return nil, err
	}
	return content, nil
}

func mapToJson(files []FileInfo) ([]byte, error) {
	resp := Message{
		Event: "server-send-files",
		Data: map[string]interface{}{
			"files": files,
		},
	}
	return json.Marshal(resp)
}
