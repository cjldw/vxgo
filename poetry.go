package vxgo

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func GetDailyPoetry() (*DailyPoetry, error) {
	request, _ := http.NewRequest(http.MethodGet, poetryURL, nil)
	request.Header.Add("X-User-Token", poetryToken)
	respData, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Printf("request daily poetry failure: %v\n", err)
		return nil, err
	}
	defer respData.Body.Close()
	respBytes, err := ioutil.ReadAll(respData.Body)
	if err != nil {
		log.Printf("read daily poetry body failure: %v\n", err)
		return nil, err
	}
	dailyPoetry := new(DailyPoetry)
	err = json.Unmarshal(respBytes, dailyPoetry)
	if err != nil {
		log.Printf("json unmarshal poetry resp failure: %v\n", err)
		return nil, err
	}
	log.Printf("get daily poetry: %#v\n", dailyPoetry)
	return dailyPoetry, nil
}
