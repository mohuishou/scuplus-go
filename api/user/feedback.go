package user

import (
	"fmt"
	"log"

	"github.com/kataras/iris"
	"github.com/mohuishou/scuplus-go/api"
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
	param := FeedBackParam{}
	err := ctx.ReadForm(&param)
	if err != nil {
		api.Error(ctx, 50002, "反馈失败", "")
		return
	}
	log.Println(param)
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
	err = github.NewIssue(param.Title, body, []string{param.Label, "用户反馈"})
	if err != nil {
		api.Error(ctx, 50001, "反馈失败", err.Error())
	}
	// 数据库保存反馈记录
	api.Success(ctx, "反馈成功", nil)
}
