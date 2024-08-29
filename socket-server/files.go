package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

type FileInfo struct {
	Path    string `json:"path"`
	Name    string `json:"name"`
	IsDir   bool   `json:"isDir"`
	Content string `json:"content,omitempty"`
}

func getFilesFromS3(userId, projectId string) {
	data := map[string]string{
		"userId":    userId,
		"projectId": projectId,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println(err.Error())
		return
	}
	url := SERVER + "/api/getUserFiles"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return
	}
	fmt.Println(string(body))
}

// TODO: after s3 try to do it for it.
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

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// TODO: Handle this more sensibly like a file can be Dockerfile
func createFileOrDir(dirpath string, files map[string]FileInfo) (map[string]FileInfo, error) {
	temp := strings.Split(dirpath, "/")
	s := temp[len(temp)-1]
	isFile := false
	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			isFile = true
		}
	}
	f := FileInfo{Path: dirpath, IsDir: !isFile, Name: s}
	files[dirpath] = f
	if !isFile {
		err := os.MkdirAll(dirpath, 0755)
		if err != nil {
			log.Fatal("error creating directory ", err)
			return nil, err
		}
		return files, nil
	}
	parentPath := ""
	for i := 0; i < len(temp)-1; i++ {
		parentPath += temp[i] + "/"
	}

	err := os.MkdirAll(parentPath, 0755)
	if err != nil {
		log.Fatal("error creating directory ", err)
		return nil, err
	}

	file, err := os.Create(dirpath)
	if err != nil {
		log.Fatal("error creating file ", err)
		return nil, err
	}
	file.Close()
	return files, nil
}

func readDir(dirpath string) (map[string]FileInfo, error) {
	files := make(map[string]FileInfo)
	//base := "/home/akhil"
	exit, err := exists(dirpath)
	if err != nil {
		log.Fatal("Error :", err)
		return nil, err
	}

	if !exit {
		f, err := createFileOrDir(dirpath, files)
		if err != nil {
			return nil, err
		}
		return f, nil
	}

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
		files[filePath] = f
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

func mapToJson(files map[string]FileInfo) ([]byte, error) {
	resp := Message{
		Event: "server-send-files",
		Data: map[string]interface{}{
			"files": files,
		},
	}
	return json.Marshal(resp)
}
