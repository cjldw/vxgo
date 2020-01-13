package vxgo

import (
	"fmt"
	"log"
	"net/http"
)

type WebSvr struct {
}

func NewWebSvr() *WebSvr {
	return new(WebSvr)
}

func (ws *WebSvr) Listen() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/blog/built.json", ws.handleBuilt)
	log.Printf("run websvr listen at: %s\n", VxCfg.WebAddr)
	return http.ListenAndServe(VxCfg.WebAddr, mux)
}

func (ws *WebSvr) handleBuilt(writer http.ResponseWriter, request *http.Request) {
	WeChatSyncRun()
	output, err := HexoDeploy()
	if err != nil {
		result := fmt.Sprintf("{\"code\": 1, \"message\":\"hexo deploy failure\", \"data\": \"%s\"}", output)
		writer.Write([]byte(result))
		return
	}
	result := fmt.Sprintf("{\"code\": 0, \"message\":\"hexo deploy success\", \"data\": \"%s\"}", output)
	writer.Write([]byte(result))
	return
}
