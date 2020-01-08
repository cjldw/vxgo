package vxgo

type VxAccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type VxNewsForm struct {
	Articles []VxNews `json:"articles"`
}

type VxNews struct {
	Title              string `json:"title"`
	ThumbMediaId       string `json:"thumb_media_id"`
	Authod             string `json:"authod"`
	ShowCoverPic       string `json:"show_cover_pic"`
	Digest             string `json:"digest"`
	Content            string `json:"content"`
	ContentSourceUrl   string `json:"content_source_url"`
	NeedOpenComment    uint32 `json:"need_open_comment"`
	OnlyFansCanComment uint32 `json:"only_fans_can_comment"`
}

type VxMaterial struct {
	URL     string `json:"url"`
	MediaId string `json:"media_id"`
}
