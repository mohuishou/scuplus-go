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
func NotifyGrade(uid uint, courseName, grade, credit string, num int) error {
	// 获取模板id
	msgID := msgid.Get(uid)
	if msgID == "" {
		return errors.New("没有模板id")
	}
	// 获取用户openid
	wechatUser := model.Wechat{}
	model.DB().Where("user_id = ?", uid).Find(&wechatUser)
	// 构造请求参数
	data := map[string]interface{}{
		"touser":      wechatUser.Openid,
		"template_id": config.Get().Wechat.TemplateGrade,
		"page":        "pages/grade",
		"form_id":     msgID,
		"data": map[string]interface{}{
			"keyword1": map[string]interface{}{
				"value": courseName,
			},
			"keyword2": map[string]interface{}{
				"value": grade,
			},
			"keyword3": map[string]interface{}{
				"value": credit,
			},
			"keyword4": map[string]interface{}{
				"value": fmt.Sprintf("共更新%d门成绩，点击查看所有成绩", num),
			},
		},
	}

	return notify(uid, data)
}

// NotifyBook 发送图书到期通知
func NotifyBook(uid uint, bookName, end string, day int64) error {
	// 获取模板id
	msgID := msgid.Get(uid)
	if msgID == "" {
		return errors.New("没有模板id")
	}
	// 获取用户openid
	wechatUser := model.Wechat{}
	model.DB().Where("user_id = ?", uid).Find(&wechatUser)
	// 构造请求参数
	data := map[string]interface{}{
		"touser":      wechatUser.Openid,
		"template_id": config.Get().Wechat.TemplateBook,
		"page":        "pages/library/loan",
		"form_id":     msgID,
		"data": map[string]interface{}{
			"keyword1": map[string]interface{}{
				"value": bookName,
			},
			"keyword2": map[string]interface{}{
				"value": end,
			},
			"keyword3": map[string]interface{}{
				"value": fmt.Sprintf("图书到期时间仅剩%d天, 点击进入我的借阅，续借图书", day),
			},
		},
	}

	return notify(uid, data)
}

// NotifyExam 发送考试提醒
func NotifyExam(uid uint, courseName, date, time, address, site, courseType string, day int64) error {
	// 获取模板id
	msgID := msgid.Get(uid)
	if msgID == "" {
		return errors.New("没有模板id")
	}
	// 获取用户openid
	wechatUser := model.Wechat{}
	model.DB().Where("user_id = ?", uid).Find(&wechatUser)
	// 构造请求参数
	data := map[string]interface{}{
		"touser":      wechatUser.Openid,
		"template_id": config.Get().Wechat.TemplateExam,
		"page":        "pages/exam",
		"form_id":     msgID,
		"data": map[string]interface{}{
			"keyword1": map[string]interface{}{
				"value": courseName,
			},
			"keyword2": map[string]interface{}{
				"value": date,
			},
			"keyword3": map[string]interface{}{
				"value": time,
			},
			"keyword4": map[string]interface{}{
				"value": address,
			},
			"keyword5": map[string]interface{}{
				"value": site,
			},
			"keyword6": map[string]interface{}{
				"value": courseType,
			},
			"keyword7": map[string]interface{}{
				"value": fmt.Sprintf("考试时间仅剩%d天，请抓紧时间复(yu)习,考试当天请携带好您的学生证，英语考试请携带听力耳机", day),
			},
		},
	}

	return notify(uid, data)
}

func NotifyFeedback(feedbackNumber int, content string) error {
	// 获取反馈信息
	feedback := model.Feedback{}
	err := model.DB().Where("number = ?", feedbackNumber).Find(&feedback).Error
	if err != nil {
		return err
	}
	uid := feedback.UserID

	// 获取模板id
	msgID := msgid.Get(uid)
	if msgID == "" {
		return errors.New("没有模板id")
	}
	// 获取用户openid
	wechatUser := model.Wechat{}
	model.DB().Where("user_id = ?", uid).Find(&wechatUser)

	// 构造请求参数
	data := map[string]interface{}{
		"touser":      wechatUser.Openid,
		"template_id": config.Get().Wechat.TemplateFeedback,
		"page":        fmt.Sprintf("pages/my/feedbackDetail?id=%d", feedback.Number),
		"form_id":     msgID,
		"data": map[string]interface{}{
			"keyword1": map[string]interface{}{
				"value": content,
			},
			"keyword2": map[string]interface{}{
				"value": feedback.Title,
			},
			"keyword3": map[string]interface{}{
				"value": feedback.Stat,
			},
			"keyword4": map[string]interface{}{
				"value": feedback.CreatedAt.Format("2006-01-02 15:04"),
			},
		},
	}
	return notify(uid, data)
}

func notify(uid uint, data map[string]interface{}) error {

	// 获取access token
	token, err := wechat.GetAccessToken(false)
	if err != nil {
		return err
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
