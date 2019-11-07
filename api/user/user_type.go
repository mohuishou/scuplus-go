package user

import (
	"github.com/kataras/iris/v12"
	"github.com/mohuishou/scuplus-go/api"
	"github.com/mohuishou/scuplus-go/middleware"
	"github.com/mohuishou/scuplus-go/model"
)

// UpdateUserType 更新用户类型
func UpdateUserType(ctx iris.Context) {
	// 获取参数
	userType, err := ctx.PostValueInt("user_type")
	if err != nil {
		api.Error(ctx, 400, "参数错误", err)
		return
	}

	// 获取用户id
	uid := middleware.GetUserID(ctx)

	// 获取用户配置
	if err := model.DB().Model(&model.UserConfig{}).Where(
		"user_id = ?",
		uid,
	).Update("user_type", userType).Error; err != nil {
		api.Error(ctx, 500, "更新失败", err)
		return
	}
	api.Success(ctx, "更新成功", nil)
}
