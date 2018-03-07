package user

import (
	"fmt"
	"time"

	"github.com/kataras/iris"
	"github.com/mohuishou/scuplus-go/api"
	"github.com/mohuishou/scuplus-go/middleware"
	"github.com/mohuishou/scuplus-go/model"
	"github.com/mohuishou/scuplus-go/util/github"
	validator "gopkg.in/go-playground/validator.v9"
)

// FeedBackParam 反馈信息的参数
type FeedBackParam struct {
	Brand          string `form:"brand" validate:"required"`          //手机品牌
	Model          string `form:"model" validate:"required"`          //手机型号
	Version        string `form:"version" validate:"required"`        //微信版本
	System         string `form:"system" validate:"required"`         //操作系统版本
	SDKVersion     string `form:"SDKVersion" validate:"required"`     //SDK版本
	ScuplusVersion string `form:"scuplusVersion" validate:"required"` //小程序版本
	Label          string `form:"label" validate:"required"`          //反馈的类型
	Title          string `form:"title" validate:"required"`          //反馈的标题
	Content        string `form:"content" validate:"required"`        //反馈的内容
}

// FeedBack 新增一条反馈信息
func FeedBack(ctx iris.Context) {

	// 频率验证，每位用户每天最多只能发送五条反馈信息
	uid := middleware.GetUserID(ctx)

	today, _ := time.ParseInLocation("2006-01-02", time.Now().Format("2006-01-02"), time.Local)
	count := 0
	model.DB().Table("feedbacks").Where("user_id = ? and created_at > ?", uid, today).Count(&count)
	if count > 5 {
		api.Error(ctx, 50004, "您今天的反馈次数已经达到上限", "")
		return
	}

	// 获取参数
	param := FeedBackParam{}
	err := ctx.ReadForm(&param)
	if err != nil {
		api.Error(ctx, 50002, "反馈失败", "")
		return
	}

	// 参数验证
	validate := validator.New()
	err = validate.Struct(param)
	if err != nil {
		api.Error(ctx, 50400, "参数错误！", err.Error())
		return
	}

	// 构造反馈内容
	body := param.Content + "\n\n\n"
	body += "scuplus_version: " + param.ScuplusVersion + "\n"
	body += fmt.Sprintf("手机: %s-%s-%s \n", param.Brand, param.Model, param.System)
	body += fmt.Sprintf("微信：%s,SDK: %s \n", param.Version, param.SDKVersion)

	// 向github新建反馈信息
	issue, err := github.NewIssue(param.Title, body, []string{param.Label, "用户反馈"})
	if err != nil {
		api.Error(ctx, 50001, "反馈失败", err.Error())
		return
	}

	// 数据库保存反馈记录
	if err := model.DB().Create(&model.Feedback{
		UserID: uid,
		Number: *issue.Number,
		Title:  param.Title,
	}).Error; err != nil {
		api.Error(ctx, 50003, "反馈失败", err.Error())
		return
	}

	api.Success(ctx, "反馈成功", nil)
}

// GetFeedBacks 获取用户的所有反馈信息
func GetFeedBacks(ctx iris.Context) {
	uid := middleware.GetUserID(ctx)
	feedbacks := []model.Feedback{}
	model.DB().Table("feedbacks").Where("user_id = ?", uid).Order("id desc").Find(&feedbacks)
	api.Success(ctx, "反馈列表获取成功！", feedbacks)
}

// GetFeedBack 获取某一条反馈信息及其所有评论
func GetFeedBack(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		api.Error(ctx, 50400, "参数错误", err)
		return
	}
	issue, comments, err := github.GetIssue(id)
	if err != nil {
		api.Error(ctx, 50002, "issue获取失败", err)
		return
	}
	api.Success(ctx, "反馈信息获取成功！", map[string]interface{}{
		"issue":    issue,
		"comments": comments,
	})
}
