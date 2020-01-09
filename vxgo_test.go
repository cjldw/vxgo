package vxgo

import (
	"fmt"
	"testing"
)

func TestGetVxNet(t *testing.T) {
	vxNET := GetVxNet()
	token := vxNET.GetAccessToken()
	t.Logf(token)
}

func TestVxNET_PostVxNews(t *testing.T) {
	vxNET := GetVxNet()

	news := VxNews{
		Title:        "欧文融通马拉碌碌无为二等分23",
		Authod:       "罗大文",
		ThumbMediaId: "drwaZ2CgYKBpJE7GXmYSXPNSX_O5SLf4P5oyx_aiMLo",
		Content:      "<h1>测试一下图片信息</h1><img src=\"https://img4.zhanqi.tv/uploads/2019/07/e9dfd0/2d7e603e7d5886f9a08d8cc32fe0a847.jpg\"/>=====",
		//ContentSourceUrl: "https://www.zhanqi.tv",
		ShowCoverPic: "1",
	}
	ok := vxNET.PostVxNews(news)
	t.Log(ok)
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
	url := vxNET.UploadVxImg("media", "E:/codelab/go/zhanqiTV/vxgo/demo.jpg")
	news := VxNews{
		Title:            "欧文融通马拉碌碌无为二等分23, 图片内问题",
		Authod:           "罗大文",
		ThumbMediaId:     "drwaZ2CgYKBpJE7GXmYSXPNSX_O5SLf4P5oyx_aiMLo",
		Content:          fmt.Sprintf("<h1>测试一下图片信息</h1><img src=\"%s\"/>=====", url),
		ContentSourceUrl: "https://loovien.github.io/2020/01/01/2020/",
		ShowCoverPic:     "1",
	}
	ok := vxNET.PostVxNews(news)
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
	PullRepo()
}

func TestGitShowCase(t *testing.T) {
	GitShowCase()
}
