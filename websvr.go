package vxgo

import (
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
	go func() {
		WeChatSyncRun()
		output, err := HexoDeploy()
		if err != nil {
			log.Printf("hexo deploy failure: %s,  %v\n", output, err)
			return
		}
		log.Printf("hexo deploy success: %s\n", output)
	}()
	writer.Write([]byte(`{"code":0, "message":"success", "data":{"tips":"job dispatched background"}}`))
	return
}
