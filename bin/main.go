package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"vxgo"
)

var (
	mode string
)

const (
	modeWeb  = "web"
	modeWx   = "wx"
	modeHexo = "hexo"
)

func init() {
	flag.StringVar(&mode, "mode", "web", "choice: [web, hexo, wx] \nweb: run web handle webhook."+
		"\nhexo: only hexo deploy.\nwx: only sync wechat.")
}

func main() {
	flag.Parse()
	switch mode {
	case modeWeb:
		sigChan := make(chan os.Signal)
		signal.Notify(sigChan, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		go func() {
			webSvr := vxgo.NewWebSvr()
			err := webSvr.Listen()
			if err != nil {
				log.Printf("websvr stop %v\n", err)
			}
		}()
		sig := <-sigChan
		log.Printf("receive terminal signal: %v, then stop websvr.\n", sig)
	case modeHexo:
		vxgo.HexoDeploy()
	case modeWx:
		vxgo.WeChatSyncRun()
	}

}
