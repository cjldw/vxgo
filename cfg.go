package vxgo

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"
)

type CFG struct {
	AppId       string `json:"appId"`
	AppSecret   string `json:"appSecret"`
	GitRepo     string `json:"gitRepo"`
	GitRepoName string `json:"gitRepoName"`
	WorkDir     string `json:"workDir"`
}

var (
	VxCfg = &CFG{}
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(file)
	cfgFile := filepath.Join(baseDir, ".vxgo")
	fileBytes, err := ioutil.ReadFile(cfgFile)
	log.Printf("load configurate file: %s\n", cfgFile)
	if err != nil {
		log.Printf("file not found: %s \n", cfgFile)
		VxCfg = &CFG{
			AppId:       appId,
			AppSecret:   appSecret,
			GitRepo:     gitRepo,
			GitRepoName: gitRepoName,
			WorkDir:     workDir,
		}
		return
	}
	err = json.Unmarshal(fileBytes, VxCfg)
	if err != nil {
		log.Fatalf("unmarshal configure file: %s failure: %v\n", cfgFile, err)
	}
}
