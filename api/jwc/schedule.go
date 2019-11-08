package jwc

import (
	"github.com/kataras/iris/v12"
	"github.com/mohuishou/scuplus-go/api"
	"github.com/mohuishou/scuplus-go/middleware"
	"github.com/mohuishou/scuplus-go/model"
)

type ScheduleParams struct {
	Term int `form:"term"`
	Year int `form:"year"`
}

// GetSchedules 获取课程表
func GetSchedules(ctx iris.Context) {
	uid := middleware.GetUserID(ctx)
	params := ScheduleParams{}
	err := ctx.ReadForm(&params)
	if err != nil {
		api.Error(ctx, 400, "参数错误", err)
		return
	}

	api.Success(ctx, "获取成功", model.GetSchedules(uid, params.Year, params.Term))
}

// UpdateSchedule 更新课程表
func UpdateSchedule(ctx iris.Context) {
	uid := middleware.GetUserID(ctx)
	params := ScheduleParams{}
	err := ctx.ReadForm(&params)
	if err != nil {
		api.Error(ctx, 400, "参数错误", err)
		return
	}
	if err := model.UpdateSchedules(uid, params.Year, params.Term); err != nil {
		api.Error(ctx, 21001, err.Error(), nil)
		return
	}
	api.Success(ctx, "更新成功", nil)
}
