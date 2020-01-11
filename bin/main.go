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
	success, _ := vxgo.GetVxNet().PostMessageBroadcast(vxm.MediaId)
	log.Printf("broadcast %s status: %v\n", vxm.MediaId, success)
}
