package jwc

import (
	"github.com/kataras/iris"
	"github.com/mohuishou/scuplus-go/api"
	"github.com/mohuishou/scuplus-go/middleware"
	"github.com/mohuishou/scuplus-go/model"
)

// GetGrade 获取成绩
func GetGrade(ctx iris.Context) {
	uid := middleware.GetUserID(ctx)
	api.Success(ctx, "获取成功", model.GetGrades(uid))
}

// UpdateGrade 更新成绩
func UpdateGrade(ctx iris.Context) {
	uid := middleware.GetUserID(ctx)

	if _, err := model.UpdateGrades(uid); err != nil {
		api.Error(ctx, 20001, "更新失败", err)
		return
	}
	api.Success(ctx, "成绩更新成功", nil)
}
