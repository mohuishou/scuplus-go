package wechat

import (
	"net/http"

	"fmt"

	"encoding/json"
	"io/ioutil"

	"github.com/kataras/iris/core/errors"
	"github.com/mohuishou/scuplus-go/cache/token"
	"github.com/mohuishou/scuplus-go/config"
)

type AccessToken struct {
	Token     string
	ExpiresIn int
}

// GetAccessToken 获取微信access token
// refresh 是否需要强制刷新，默认优先获取缓存数据
func GetAccessToken(refresh bool) (string, error) {
	// 获取缓存token
	t := token.Get()
	if !refresh && t != "" {
		return t, nil
	}

	// 从微信服务器获取token
	c := http.Client{}
	resp, err := c.Get(fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		config.Get().Wechat.Appid,
		config.Get().Wechat.Secret))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	at := AccessToken{}
	json.Unmarshal(body, &at)

	// 缓存token并且返回
	if at.Token == "" {
		return "", errors.New("获取token失败！")
	}
	err = token.Set(at.Token)
	if err != nil {
		return "", err
	}
	return at.Token, nil
}
