package user

import (
	"github.com/kataras/iris/v12"
	"github.com/mohuishou/scuplus-go/api"
	"github.com/mohuishou/scuplus-go/middleware"
	"github.com/mohuishou/scuplus-go/model"
)

func UpdateNotify(ctx iris.Context) {
	notify, err := ctx.PostValueInt("notify")
	if err != nil {
		api.Error(ctx, 400, "参数错误", err)
		return
	}

	uid := middleware.GetUserID(ctx)
	if err := model.DB().Model(&model.UserConfig{}).Where("user_id = ?", uid).Update("notify", notify).Error; err != nil {
		api.Error(ctx, 500, "更新失败", err)
		return
	}
	api.Success(ctx, "更新成功！", nil)
}

func GetNotify(ctx iris.Context) {
	uid := middleware.GetUserID(ctx)
	userConf := model.UserConfig{}
	model.DB().Where("user_id = ?", uid).Select([]string{"notify"}).Find(&userConf)
	api.Success(ctx, "获取成功", userConf.Notify)
}
