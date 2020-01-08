package vxgo

import (
	"io/ioutil"
	"log"
)

func GetDirFiles(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalf("read directory failure: %v\n", err)
	}
	fileLen := len(files)
	var respData []string
	for i := 0; i < fileLen; i++ {
		for j := 1; j < fileLen-1; j++ {
			prevFile := files[i]
			nextFile := files[j]
			if prevFile.ModTime().Before(nextFile.ModTime()) {
				files[i], files[j] = nextFile, prevFile
			}
		}
	}
	for i := 0; i < fileLen; i++ {
		if files[i].IsDir() {
			continue
		}
		respData = append(respData, files[i].Name())
	}
	return respData
}
