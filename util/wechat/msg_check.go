package wechat

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/json-iterator/go"
)

// MsgCheck 文字安全检查
func MsgCheck(content string) (bool, error) {
	c := http.Client{}
	token, err := GetAccessToken(false)
	if err != nil {
		return false, err
	}
	param := map[string]string{
		"content": content,
	}
	paramB, _ := jsoniter.Marshal(&param)

	resp, err := c.Post(
		"https://api.weixin.qq.com/wxa/msg_sec_check?access_token="+token,
		"application/json",
		bytes.NewReader(paramB),
	)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	data := struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}{}
	err = jsoniter.Unmarshal(body, &data)
	if err != nil {
		return false, err
	}
	if data.ErrCode != 0 {
		return false, err
	}
	return true, nil
}
