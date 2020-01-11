package vxgo

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var weChatTPL = `
{{ if poetry }}
    <h1>{{.poetry.Data.Origin.Title}}</h1>
	{{ range .poetry.Data.Origin.Content}}
		<p>{{.element}}</p>
	{{end}}
	<span>{{poetry.Data.Origin.Author}}</span>
{{else}}
    <h1>今日诗歌歇菜了</h1>
	<p>................</p>
{{end}}

<p> 点击左下角阅读更多</p>
`

func ParseVxNews(file string) (*VxNews, error) {

	fd, err := os.Open(file)
	if err != nil {
		log.Printf("open file: %s failure: %v\n", file, err)
		return nil, err
	}
	scanner := bufio.NewScanner(fd)
	poetry, _ := GetDailyPoetry()
	writer := bytes.NewBufferString(weChatTPL)
	_ = template.New("WeChatTPL").Execute(writer, poetry)

	vxNews := &VxNews{
		ThumbMediaId: "drwaZ2CgYKBpJE7GXmYSXPNSX_O5SLf4P5oyx_aiMLo",
		Content:      writer.String(),
		ShowCoverPic: "1",
	}

	for scanner.Scan() {
		txt := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(txt, "<!--more-->") {
			break
		}
		if strings.HasPrefix(txt, "title") {
			title := strings.Split(txt, ":")
			if len(title) >= 2 {
				vxNews.Title = strings.TrimSpace(title[1])
			}
			continue
		}
		if strings.HasPrefix(txt, "date") {
			date := strings.Split(txt, "date:")
			if len(date) < 2 {
				vxNews.ContentSourceUrl = blogURL
				continue
			}
			blogTime, err := time.Parse("2006-01-02 15:04:05", strings.TrimSpace(date[1]))
			if err != nil {
				continue
			}
			baseName := strings.TrimRight(filepath.Base(file), ".md")
			vxNews.ContentSourceUrl = fmt.Sprintf(
				"%s/%d/%s/%s/%s",
				blogURL,
				blogTime.Year(),
				blogTime.Format("01"),
				blogTime.Format("02"),
				baseName,
			)
			continue
		}
		if strings.HasPrefix(txt, "desc") {
			desc := strings.Split(txt, ":")
			if len(desc) >= 2 {
				vxNews.Digest = strings.TrimSpace(desc[1])
			}
			continue
		}
	}
	return vxNews, nil

}
