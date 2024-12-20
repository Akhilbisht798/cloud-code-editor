package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/websocket"
)

type FileInfo struct {
	Path    string `json:"path"`
	Name    string `json:"name"`
	IsDir   bool   `json:"isDir"`
	Content string `json:"content,omitempty"`
}

func getFilesFromS3(userId, projectId string) error {
	log.Println("getting files from s3")
	data := map[string]string{
		"userId":    userId,
		"projectId": projectId,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	url := SERVER + "/api/getUserFiles"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err.Error())
		return err
	}
	log.Println("Fetched presigned urls.")
	defer resp.Body.Close()

	urls := make(map[string]string)
	err = json.NewDecoder(resp.Body).Decode(&urls)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	log.Println("Length of files: ", urls)

	for k, v := range urls {
		fmt.Println("Getting file ", v)
		err = saveFileFromS3(v, k)
		if err != nil {
			log.Println("Error saving the file: ", err.Error())
			return err
		}
	}
	log.Println("File process done.")
	return nil
}

func saveFileFromS3(path string, url string) error {
	log.Println("Saving file: ", path)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println("BODY: ", string(body))
	err = saveFile(path, string(body))
	if err != nil {
		return err
	}

	return nil
}

func saveFile(path string, content string) error {
	log.Println("Saving File Locally: ", path)
	if err := os.MkdirAll(filepath.Dir(path), 0770); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(content))
	if err != nil {
		log.Println("Error writing the file: ", err)
		return err
	}
	fmt.Println("File Saved successfully.")
	return nil
}

func saveDir(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}
	return nil
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

func newFileOrDir(conn *websocket.Conn, msg Message) {
	fileData, ok := msg.Data["file"].(map[string]interface{})
	if !ok {
		fmt.Println("error ", ok)
		return
	}
	filepath := fileData["path"].(string) + "/" + fileData["name"].(string)
	files := make(map[string]FileInfo)
	var file FileInfo
	if fileData["isDir"] == false {
		err := saveFile(filepath, "")
		if err != nil {
			fmt.Println(err)
			return
		}
		file = FileInfo{
			Path:    fileData["path"].(string),
			IsDir:   false,
			Content: "",
			Name:    fileData["name"].(string),
		}
	} else {
		err := saveDir(filepath)
		if err != nil {
			fmt.Println(err)
			return
		}
		file = FileInfo{
			Path:  fileData["path"].(string),
			IsDir: true,
			Name:  fileData["name"].(string),
		}
	}

	files[filepath] = file
	resp, err := mapToJson(files)
	if err != nil {
		fmt.Println(err)
		return
	}
	conn.WriteMessage(websocket.TextMessage, resp)
}

func deleteFileOrDir(msg Message) {
	fileData, ok := msg.Data["file"].(map[string]interface{})
	if !ok {
		fmt.Println("error ", ok)
		return
	}
	filepath := fileData["path"].(string) + "/" + fileData["name"].(string)
	err := os.RemoveAll(filepath)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
}
