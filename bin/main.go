package main

import (
	"log"
	"vxgo"
)

func main() {
	_ = vxgo.PullRepo()
	commit, err := vxgo.GitShowCase()
	if err != nil {
		log.Fatalf("git commit log failure: %v\n", err)
	}
	_, exists, _ := vxgo.GetDumper().QueryCommit(commit.CommitID)
	if exists {
		log.Fatalf("this commit id processed")
	}
	var vxNewsList []*vxgo.VxNews
	for i := 0; i < len(commit.Files); i++ {
		file := commit.Files[i]
		vxNews, err := vxgo.ParseVxNews(file)
		if err != nil {
			log.Printf("parse WeChat News failure: %v\n", err)
			continue
		}
		vxNewsList = append(vxNewsList, vxNews)
	}
	vxm, err := vxgo.GetVxNet().PostVxNews(vxNewsList)
	if err != nil {
		log.Fatalf("post WeChat News failure: %v\n", err)
	}
	log.Printf("post WeChat news status: %v\n", vxm)
	// personal WeChat Account no this privileges
	// success, _ := vxgo.GetVxNet().PostNewsBroadcast(vxm.MediaId, nil)
	// log.Printf("broadcast %s status: %v\n", vxm.MediaId, success)
	success, err := vxgo.GetDumper().SaveCommit(commit.CommitID, true)
	if err != nil {
		log.Fatalf("save commit:%s failure: %v\n", commit.CommitID, err)
	}
	log.Printf("save commit: %s status: %v\n", commit.CommitID, success)
}
