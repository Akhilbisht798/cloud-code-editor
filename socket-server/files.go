package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

type FileInfo struct {
	Path    string `json:"path"`
	Name    string `json:"name"`
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
	f, err := mapToJson(files, path)
	if err != nil {
		log.Println(err)
		return
	}
	conn.WriteMessage(websocket.TextMessage, []byte(f))
}

func readDir(dirpath string) ([]FileInfo, error) {
	files := make([]FileInfo, 0)
	//base := "/home/akhil"

	entries, err := os.ReadDir(dirpath)
	if err != nil {
		log.Fatal("error reading directory: ", err)
		return nil, err
	}
	for _, e := range entries {
		filePath := dirpath + "/" + e.Name()
		f := FileInfo{Path: dirpath, IsDir: e.IsDir(), Name: e.Name()}
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

func mapToJson(files []FileInfo, path string) ([]byte, error) {
	resp := Message{
		Event: "server-send-files",
		Data: map[string]interface{}{
			"files": files,
			"dir":   path,
		},
	}
	return json.Marshal(resp)
}
