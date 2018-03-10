package jwc

import (
	"github.com/kataras/iris"
	"github.com/mohuishou/scuplus-go/api"
	"github.com/mohuishou/scuplus-go/middleware"
	"github.com/mohuishou/scuplus-go/model"
)

// GetSchedules 获取课程表
func GetSchedules(ctx iris.Context) {
	uid := middleware.GetUserID(ctx)
	term := ctx.FormValue("term")
	api.Success(ctx, "获取成功", model.GetSchedules(uid, term))

}

// UpdateSchedule 更新课程表
func UpdateSchedule(ctx iris.Context) {
	uid := middleware.GetUserID(ctx)
	term := ctx.FormValue("term")
	if err := model.UpdateSchedules(uid, term); err != nil {
		api.Error(ctx, 21001, err.Error(), nil)
		return
	}
	api.Success(ctx, "更新成功", nil)
}
