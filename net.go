package vxgo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type VxNET struct {
	accessToken *VxAccessToken
}

var (
	vxNET     *VxNET
	VxNetOnce = &sync.Once{}
)

func GetVxNet() *VxNET {
	VxNetOnce.Do(func() {
		vxNET = &VxNET{
			accessToken: &VxAccessToken{},
		}
	})
	return vxNET
}

func (vn *VxNET) GetAccessToken() (string, error) {
	if vn.accessToken.ExpiresIn > int(time.Now().Unix()) {
		return vn.accessToken.AccessToken, nil
	}
	tokenURL := fmt.Sprintf(accessTokenURL, VxCfg.AppId, VxCfg.AppSecret)
	response, err := http.Get(tokenURL)
	if err != nil {
		log.Printf("fetch WeChat access_token failure:%v\n", err)
		return "", err
	}
	defer response.Body.Close()
	respBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("read access_token response failure: %#v\n", err)
		return "", err
	}
	accessToken := &VxAccessToken{}
	err = json.Unmarshal(respBytes, accessToken)
	if err != nil {
		log.Printf("unmarshal access_token failure: %#v\n", err)
		return "", err
	}
	if len(accessToken.AccessToken) <= 0 {
		log.Printf("get WeChat access_token failure: %s\n", string(respBytes))
		return "", errors.New(string(respBytes))
	}
	vn.accessToken = &VxAccessToken{
		AccessToken: accessToken.AccessToken,
		ExpiresIn:   accessToken.ExpiresIn + int(time.Now().Unix()),
	}
	return accessToken.AccessToken, nil
}

func (vn *VxNET) PostPersistentMaterial(fileName, filePath, typ string, params ...map[string]string) *VxMaterial {
	accessToken, err := vn.GetAccessToken()
	if err != nil {
		log.Fatal(err)
	}
	mURL := fmt.Sprintf(materialURL, accessToken, typ)
	buffer := &bytes.Buffer{}
	partBody := multipart.NewWriter(buffer)

	part, err := partBody.CreateFormFile(fileName, filepath.Base(filePath))
	if err != nil {
		log.Fatalf("create multipart filePath failure: %#v\n", err)
	}
	fd, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("post material to WeChat failure: %s err: %#v\n", filePath, err)
	}
	io.Copy(part, fd)
	fd.Close()
	if len(params) > 0 {
		for key, value := range params[0] {
			err = partBody.WriteField(key, value)
			if err != nil {
				log.Printf("write field: %s value: %s failure: %#v\n", key, value, err)
			}
		}
	}
	contentType := partBody.FormDataContentType()
	resp, err := http.Post(mURL, contentType, buffer)
	partBody.Close()
	if err != nil {
		log.Fatalf("send material request failure: %#v\n", err)
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("read material request failure: %s err: %#v", string(respBytes), err)
	}
	material := &VxMaterial{}
	log.Printf("send material request success: %s\n", string(respBytes))
	err = json.Unmarshal(respBytes, material)
	if err != nil {
		log.Fatalf("post material response unmarshal failure: %v\n", err)
	}
	return material
}

func (vn *VxNET) PostVxNews(news []*VxNews) (*VxMaterial, error) {
	token, err := vn.GetAccessToken()
	if err != nil {
		return nil, err
	}
	newsUrl := fmt.Sprintf(addNewsURL, token)
	vxNews := VxNewsForm{
		Articles: news,
	}
	newBytes, _ := json.Marshal(vxNews)
	formReader := bytes.NewReader(newBytes)
	response, err := http.Post(newsUrl, "application/json;charset=utf8", formReader)
	if err != nil {
		log.Printf("post news to WeChat failure: data:%#v,  err: %#v\n", vxNews, err)
		return nil, err
	}
	defer response.Body.Close()
	respBytes, _ := ioutil.ReadAll(response.Body)
	log.Printf("post news to WeChat result: %s\n", string(respBytes))
	vxMaterial := new(VxMaterial)
	err = json.Unmarshal(respBytes, vxMaterial)
	if err != nil {
		log.Printf("unmarshal post news response failure: %v\n", err)
		return nil, err
	}
	return vxMaterial, nil
}

func (vn *VxNET) UploadVxImg(paramName, file string) (string, error) {
	fd, err := os.Open(file)
	if err != nil {
		log.Fatalf("open file: %s failure: %v\n", file, err)
	}
	body := new(bytes.Buffer)
	part := multipart.NewWriter(body)
	partBody, err := part.CreateFormFile(paramName, filepath.Base(file))
	if err != nil {
		log.Printf("create form upload file failure: %v\n", err)
		return "", err
	}
	io.Copy(partBody, fd)
	fd.Close()
	part.Close()

	token, err := vn.GetAccessToken()
	if err != nil {
		return "", err
	}
	uploadURL := fmt.Sprintf(imgUploadingURL, token)
	resp, err := http.Post(uploadURL, part.FormDataContentType(), body)
	if err != nil {
		log.Printf("request uploadimg: %s failure: %v\n", uploadURL, err)
		return "", err
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("read response failure: %v\n", err)
		return "", err
	}
	resp.Body.Close()
	log.Printf("request uploadimg success: %s\n", string(respBytes))

	type ImgUploadResp struct {
		URL string `json:"url"`
	}
	imgResp := new(ImgUploadResp)
	err = json.Unmarshal(respBytes, imgResp)
	if err != nil {
		log.Printf("json unmarshal failure: %v\n", err)
		return "", err
	}
	return imgResp.URL, nil
}

func (vn *VxNET) PostNewsBroadcast(mediaId string, filter *MessSendFilter) (bool, error) {
	token, err := vn.GetAccessToken()
	if err != nil {
		return false, err
	}
	broadcastURL := fmt.Sprintf(messPushURL, token)
	newsSend := NewsSend{
		MsgType:           "mpnews",
		Filter:            filter,
		MpNews:            MediaId{MediaId: mediaId},
		SendIgnoreReprint: 0,
	}
	postBytes, err := json.Marshal(newsSend)
	if err != nil {
		log.Printf("json marshal news post data failure: %v\n", err)
		return false, err
	}
	resp, err := http.Post(broadcastURL, "application/json;charset=utf-8", bytes.NewBuffer(postBytes))
	if err != nil {
		log.Printf("post news to WeChat broadcast failure: %v\n", err)
		return false, err
	}
	defer resp.Body.Close()

	respBytes, _ := ioutil.ReadAll(resp.Body)
	log.Printf("post news to WeChat broadcast status: %s\n", string(respBytes))
	return true, nil
}

func (vn *VxNET) MessageBroadcast(mediaId, mediaType string) {

}
