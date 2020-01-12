package vxgo

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"
	"time"
)

func TestGetVxNet(t *testing.T) {
	vxNET := GetVxNet()
	token, err := vxNET.GetAccessToken()
	t.Logf(token, err)
}

func TestVxNET_PostVxNews(t *testing.T) {
	vxNET := GetVxNet()

	news := &VxNews{
		Title:        "欧文融通马拉碌碌无为二等分23",
		Authod:       "罗大文",
		ThumbMediaId: "drwaZ2CgYKBpJE7GXmYSXPNSX_O5SLf4P5oyx_aiMLo",
		Content:      "<h1>测试一下图片信息</h1><img src=\"https://img4.zhanqi.tv/uploads/2019/07/e9dfd0/2d7e603e7d5886f9a08d8cc32fe0a847.jpg\"/>=====",
		//ContentSourceUrl: "https://www.zhanqi.tv",
		ShowCoverPic: "1",
	}
	ok, _ := vxNET.PostVxNews([]*VxNews{news})
	t.Log(ok)
}

func TestVxNET_GetAccessToken(t *testing.T) {
	token, err := GetVxNet().GetAccessToken()
	t.Log(token, err)
}

func TestVxNET_PostPersistentMaterial(t *testing.T) {
	vxNET := GetVxNet()
	params := map[string]string{
		"description": `{"title":一个测试图片素材, "introduction":"哈哈, 牛逼的, 一起来吧吧"}`,
	}
	rs := vxNET.PostPersistentMaterial("media", "E:/codelab/go/zhanqiTV/vxgo/demo.jpg", imageType, params)
	t.Log(rs.MediaId, rs.URL)
}

func TestVxNET_UploadVxImg(t *testing.T) {
	vxNET := GetVxNet()
	url, _ := vxNET.UploadVxImg("media", "E:/codelab/go/zhanqiTV/vxgo/demo.jpg")
	news := &VxNews{
		Title:            "欧文融通马拉碌碌无为二等分23, 图片内问题",
		Authod:           "罗大文",
		ThumbMediaId:     "drwaZ2CgYKBpJE7GXmYSXPNSX_O5SLf4P5oyx_aiMLo",
		Content:          fmt.Sprintf("<h1>测试一下图片信息</h1><img src=\"%s\"/>=====", url),
		ContentSourceUrl: "https://loovien.github.io/2020/01/01/2020/",
		ShowCoverPic:     "1",
	}
	ok, _ := vxNET.PostVxNews([]*VxNews{news})
	t.Log(ok)
}

func TestGetDirFiles(t *testing.T) {
	files := GetDirFiles("D:/codelab/zhanqiTV/go/vxgo")
	t.Logf("%#v\n", files)
}

func TestGitClone(t *testing.T) {
	CloneRepo()
}

func TestPullRepo(t *testing.T) {
	err := PullRepo()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGitShowCase(t *testing.T) {
	commit, err := GitShowCase()
	t.Log(commit, err)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetDailyPoetry(t *testing.T) {
	r, err := GetDailyPoetry()
	t.Log(r, err)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDemo(t *testing.T) {
	t.Log(filepath.Base("/luowen/haha/sdb/c/we/abc.md"))
	blogTime := time.Now()
	a := fmt.Sprintf(
		"%s/%d/%s/%d/%s",
		blogURL,
		blogTime.Year(),
		blogTime.Month(),
		blogTime.Day(),
		"文集案例大家",
	)
	t.Log(a)

}

func TestParseVxNews(t *testing.T) {
	vx, e := ParseVxNews("/tmp/blog-code/source/_posts/2020.md")
	t.Log(vx, e)
	if e != nil {
		t.Fatal(e)
	}
}

func TestParsePoetry(t *testing.T) {
	vxm := VxMaterial{
		URL:     "http://www.qq.com",
		MediaId: "aaabbbcccee",
	}
	b, e := json.Marshal(vxm)
	t.Log(vxm, string(b), e)
	//c := parsePoetry()
	//t.Log(c)
}

func TestGetDumper(t *testing.T) {
	dump := GetDumper()
	b, e := dump.SaveCommit("isLoad", true)
	t.Log(b, e)
}
