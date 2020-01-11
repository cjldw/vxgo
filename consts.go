package vxgo

const (
	appId       = "wx59bce618019e60cf"
	appSecret   = "7c6b8bc4ecb1c78aba481c52014f572f"
	gitRepo     = "https://github.com/loovien/blog-code.git"
	blogURL     = "https://loovien.github.io"
	gitRepoName = "blog-code"
	workDir     = "/tmp"
)

const (
	imageType = "image"
	voiceType = "voice"
	videoType = "video"
	thumbType = "thumb"
)

const (
	addNewsURL      = "https://api.weixin.qq.com/cgi-bin/material/add_news?access_token=%s"
	accessTokenURL  = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	materialURL     = "https://api.weixin.qq.com/cgi-bin/material/add_material?access_token=%s&type=%s"
	imgUploadingURL = "https://api.weixin.qq.com/cgi-bin/media/uploadimg?access_token=%s"
	messPushURL     = "https://api.weixin.qq.com/cgi-bin/message/mass/sendall?access_token=%s"
)

const (
	poetryToken = "elS8gi/k+rfUQRJI67sEgPF/mUuKDzqU"
	poetryURL   = "https://v2.jinrishici.com/sentence"
)
