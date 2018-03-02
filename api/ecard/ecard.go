package ecard

import (
	"github.com/kataras/iris"
	"github.com/mohuishou/scuplus-go/api"
	"github.com/mohuishou/scuplus-go/middleware"
	"github.com/mohuishou/scuplus-go/model"
)

// Params 参数
type Params struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

// Get 获取最新的一卡通信息
func Get(ctx iris.Context) {
	uid := middleware.GetUserID(ctx)
	params := Params{}
	if err := ctx.ReadForm(&params); err != nil {
		api.Error(ctx, 60400, "参数错误", nil)
		return
	}
	var ecards []model.Ecard
	model.DB().Where("user_id = ?", uid).Order("trans_time desc").Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize).Find(&ecards)

	api.Success(ctx, "一卡通信息获取成功", ecards)
}

// Update 更新一卡通信息
func Update(ctx iris.Context) {
	uid := middleware.GetUserID(ctx)
	if err := model.UpdateEcard(uid); err != nil {
		api.Error(ctx, 60001, "更新失败", err.Error())
		return
	}
	api.Success(ctx, "更新成功", nil)
}
