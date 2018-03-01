package tasks

import (
	"errors"
	"net/http"

	"fmt"

	"encoding/json"

	"bytes"

	"io/ioutil"

	"github.com/mohuishou/scuplus-go/cache/msgid"
	"github.com/mohuishou/scuplus-go/config"
	"github.com/mohuishou/scuplus-go/model"
	"github.com/mohuishou/scuplus-go/util/wechat"
)

// NotifyGrade 发送成绩更新通知
func NotifyGrade(uid uint, grade model.Grade, num int) error {
	// 获取模板id
	msgID := msgid.Get(uid)
	if msgID == "" {
		return errors.New("没有模板id")
	}

	// 获取access token
	token, err := wechat.GetAccessToken(false)
	if err != nil {
		return err
	}

	// 获取用户openid
	wechatUser := model.Wechat{}
	model.DB().Where("user_id = ?", uid).Find(&wechatUser)
	// 构造请求参数
	data := map[string]interface{}{
		"touser":      wechatUser.Openid,
		"template_id": config.Get().Wechat.TemplateGrade,
		"page":        "grade",
		"form_id":     msgID,
		"data": map[string]interface{}{
			"keyword1": map[string]interface{}{
				"value": grade.CourseName,
			},
			"keyword2": map[string]interface{}{
				"value": grade.Grade,
			},
			"keyword3": map[string]interface{}{
				"value": grade.Credit,
			},
			"keyword4": map[string]interface{}{
				"value": fmt.Sprintf("共更新%d门成绩，点击查看所有成绩", num),
			},
		},
	}
	b, err := json.Marshal(&data)
	if err != nil {
		return err
	}
	body := bytes.NewReader(b)

	// 发送模板消息
	c := http.Client{}
	resp, err := c.Post(
		"https://api.weixin.qq.com/cgi-bin/message/wxopen/template/send?access_token="+token,
		"",
		body,
	)
	if err != nil {
		return err
	}
	// 解析返回数据，判断是否发送成功
	defer resp.Body.Close()
	var res map[string]interface{}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	json.Unmarshal(resBody, &res)
	if code, ok := res["errcode"].(int); !ok || code != 0 {
		return errors.New(res["errmsg"].(string))
	}
	return nil
}
