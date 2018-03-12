package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/kataras/iris"
	"github.com/mohuishou/scuplus-go/api"
	"github.com/mohuishou/scuplus-go/util/wechat"
)

type QCodeParams struct {
	Page      string `form:"page" json:"page"`
	Scene     string `form:"scene" json:"scene"`
	AutoColor bool   `json:"auto_color"`
}

// GetQCode 获取二维码图片
func GetQCode(ctx iris.Context) {
	// 获取参数
	var params QCodeParams
	params.AutoColor = false
	params.Page = ctx.URLParam("page")
	params.Scene = ctx.URLParam("scene")
	if params.Page == "" || params.Scene == "" {
		api.Error(ctx, 70400, "参数错误！", nil)
		return
	}
	b, err := json.Marshal(&params)
	if err != nil {
		api.Error(ctx, 70400, "参数错误！", err)
		return
	}
	log.Println(string(b))
	// 获取token
	t, err := wechat.GetAccessToken(false)
	if err != nil {
		api.Error(ctx, 70400, "系统错误", err)
		return
	}
	// 建立http客户端
	c := http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=%s", t), bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	if err != nil {
		api.Error(ctx, 70001, "数据获取错误！", err)
		return
	}
	resp, err := c.Do(req)
	// resp, err := c.Post(fmt.Sprintf("https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=%s", t), "application/json", strings.NewReader("page=pages/details&sence=1"))
	if err != nil {
		api.Error(ctx, 70001, "数据获取错误！", err)
		return
	}
	defer resp.Body.Close()
	ctx.ContentType(resp.Header.Get("Content-Type"))
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		api.Error(ctx, 70001, "数据获取错误！", err)
		return
	}
	ctx.Write(b)
}
