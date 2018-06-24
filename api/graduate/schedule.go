package graduate

import (
	"github.com/kataras/iris"
	"github.com/mohuishou/scuplus-go/api"
	"github.com/mohuishou/scuplus-go/middleware"
	"github.com/mohuishou/scuplus-go/model"
)

type ScheduleParams struct {
	Term int `form:"term"`
	Year int `form:"year"`
}

func GetSchedule(ctx iris.Context) {
	uid := middleware.GetUserID(ctx)
	params := ScheduleParams{}
	err := ctx.ReadForm(&params)
	if err != nil {
		api.Error(ctx, 400, "参数错误", err)
		return
	}
	schdules := []model.GraduateSchedule{}
	model.DB().Where(
		"user_id = ? and year = ? and term = ?",
		uid,
		params.Year,
		params.Term,
	).Find(&schdules)
	api.Success(ctx, "获取成功", schdules)
}

func UpdateSchedule(ctx iris.Context) {
	uid := middleware.GetUserID(ctx)
	params := ScheduleParams{}
	err := ctx.ReadForm(&params)
	if err != nil {
		api.Error(ctx, 400, "参数错误", err)
		return
	}
	if err := model.UpdateGraduateSchedule(uid, params.Year, params.Term); err != nil {
		api.Error(ctx, 500, err.Error(), nil)
		return
	}
	api.Success(ctx, "更新成功", nil)
}
