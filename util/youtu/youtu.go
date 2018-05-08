// Package youtu 腾讯优图相关操作
// 这里实现app鉴权以及名片识别
package youtu

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/mohuishou/scuplus-go/config"

	"github.com/json-iterator/go"
)

const expiredInterval = 1000

func orignalSign() string {
	now := time.Now().Unix()
	rand.Seed(int64(now))
	rnd := rand.Int31()
	sign := fmt.Sprintf("u=%s&a=%s&k=%s&e=%d&t=%d&r=%d&f=",
		config.Get().Youtu.QQ,
		config.Get().Youtu.AppID,
		config.Get().Youtu.SecretID,
		now+expiredInterval,
		now,
		rnd,
	)
	return sign
}

func sign() string {
	origSign := orignalSign()
	h := hmac.New(sha1.New, []byte(config.Get().Youtu.SecretKey))
	h.Write([]byte(origSign))
	hm := h.Sum(nil)
	//attach orig_sign to hm
	dstSign := []byte(string(hm) + origSign)
	b64 := base64.StdEncoding.EncodeToString(dstSign)
	return b64
}

func post(addr, param string) (rsp []byte, err error) {
	client := &http.Client{
		Timeout: time.Duration(15 * time.Second),
	}
	httpreq, err := http.NewRequest("POST", addr, strings.NewReader(param))
	if err != nil {
		return
	}

	httpreq.Header.Add("Authorization", sign())
	httpreq.Header.Add("Content-Type", "text/json")
	httpreq.Header.Add("User-Agent", "")
	httpreq.Header.Add("Accept", "*/*")
	httpreq.Header.Add("Expect", "100-continue")
	resp, err := client.Do(httpreq)

	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		errStr := fmt.Sprintf("httperrorcode: %d \n", resp.StatusCode)
		err = errors.New(errStr)
		return
	}

	defer resp.Body.Close()
	rsp, err = ioutil.ReadAll(resp.Body)
	return
}

type OCRResp struct {
	ErrorCode int    `json:"errorcode"`
	ErrorMsg  string `json:"errormsg"`
	Items     []struct {
		Item       string `json:"item"`
		ItemString string `json:"itemstring"`
	} `json:"items"`
}

// OCR 识别传递过来的链接
func OCR(url string) (*OCRResp, error) {
	param := map[string]string{
		"url":    url,
		"app_id": config.Get().Youtu.AppID,
	}
	paramStr, err := jsoniter.MarshalToString(&param)
	if err != nil {
		return nil, err
	}
	b, err := post("http://api.youtu.qq.com/youtu/ocrapi/bcocr", paramStr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var resp OCRResp
	jsoniter.Unmarshal(b, &resp)
	return &resp, nil
}
