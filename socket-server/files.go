package main

import (
	"encoding/json"
	"log"
	"os"
)

type FileInfo struct {
	Path    string `json:"path"`
	IsDir   bool   `json:"isDir"`
	Content string `json:"content,omitempty"`
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
		f := FileInfo{Path: filePath, IsDir: e.IsDir()}
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
	return json.Marshal(files)
}
