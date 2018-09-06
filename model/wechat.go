package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/mohuishou/scuplus-go/config"
)

// Wechat 用户微信相关的信息
type Wechat struct {
	Model
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	NickName   string `json:"nick_name"`
	AvatarURL  string `json:"avatar_url"`
	Gender     string `json:"gender"`
	UserID     uint   `json:"user_id"`
}

// GetOpenid 从微信服务器获取openid等信息
func (w *Wechat) GetOpenid(code string) error {
	c := http.Client{}
	resp, err := c.Get(fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", config.Get().Wechat.Appid, config.Get().Wechat.Secret, code))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println(string(body))
	json.Unmarshal(body, &w)
	return nil
}
