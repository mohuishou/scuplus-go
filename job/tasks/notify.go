package tasks

import (
	"errors"
	"net/http"

	"fmt"

	"encoding/json"

	"bytes"

	"io/ioutil"

	"log"

	"github.com/mohuishou/scuplus-go/cache/msgid"
	"github.com/mohuishou/scuplus-go/config"
	"github.com/mohuishou/scuplus-go/model"
	"github.com/mohuishou/scuplus-go/util/wechat"
)

// NotifyGrade 发送成绩更新通知
func NotifyGrade(uid uint, courseName string, grade float64, credit string, num int) error {
	if !notifyCheck(uid, "grade") {
		log.Println("用户关闭了成绩通知", uid)
		return nil
	}
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
	if !notifyCheck(uid, "library") {
		log.Println("用户关闭了图书通知", uid)
		return nil
	}

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
	if !notifyCheck(uid, "exam") {
		log.Println("用户关闭了考试通知", uid)
		return nil
	}

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

// NotifyLostFind 给失主发送消息
// userID 用户id
// id 失物招领id
func NotifyLostFind(uid, id uint) error {
	// 获取失物招领详情
	lostFind := model.LostFind{}
	if err := model.DB().Find(&lostFind, id).Error; err != nil {
		return err
	}

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
		"template_id": config.Get().Wechat.TemplateLostFind,
		"page":        fmt.Sprintf("pages/lostFind/item?id=%d", id),
		"form_id":     msgID,
		"data": map[string]interface{}{
			"keyword1": map[string]interface{}{
				"value": lostFind.Title,
			},
			"keyword2": map[string]interface{}{
				"value": lostFind.Info,
			},
			"keyword3": map[string]interface{}{
				"value": lostFind.CreatedAt.Format("2006-01-02 15:04"),
			},
			"keyword4": map[string]interface{}{
				"value": lostFind.Contact,
			},
			"keyword5": map[string]interface{}{
				"value": "您的校园卡被同学拾到了，点击查看详情",
			},
		},
	}
	return notify(uid, data)
}

// NotifyLostFind 给失主发送消息
// id 失物招领id
// progress 进度
// comment 备注
func NotifyLostFindStatus(id uint, progress, comment string) error {
	// 获取失物招领详情
	lostFind := model.LostFind{}
	if err := model.DB().Find(&lostFind, id).Error; err != nil {
		return err
	}

	uid := lostFind.UserID
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
		"template_id": config.Get().Wechat.TemplateLostFindStatus,
		"page":        fmt.Sprintf("pages/lostFind/item?id=%d", id),
		"form_id":     msgID,
		"data": map[string]interface{}{
			"keyword1": map[string]interface{}{
				"value": progress,
			},
			"keyword2": map[string]interface{}{
				"value": comment,
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

func notifyCheck(uid uint, notifyType string) bool {
	userConf := model.UserConfig{}
	if err := model.DB().Where("user_id = ?", uid).Find(&userConf).Error; err != nil {
		log.Println("用户配置获取错误：", uid, err)
		return true
	}
	switch notifyType {
	case "grade":
		return (userConf.Notify & model.NotifyGrade) == model.NotifyGrade
	case "exam":
		return (userConf.Notify & model.NotifyExam) == model.NotifyExam
	case "library":
		return (userConf.Notify & model.NotifyLibrary) == model.NotifyLibrary
	}
	return true
}
