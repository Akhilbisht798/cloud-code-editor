package main

import (
	"log"
	"os"
)

func readDir(dirpath string) (map[string]string, error) {
	files := make(map[string]string)
	entries, err := os.ReadDir(dirpath)
	if err != nil {
		log.Fatal("error reading directory: ", err)
		return nil, err
	}
	for _, e := range entries {
		filePath := dirpath + "/" + e.Name()
		if !e.IsDir() {
			fileContent, err := readFile(filePath)
			if err != nil {
				return nil, err
			}
			files[filePath] = string(fileContent)
		} else {
			files[filePath] = ""
		}
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


