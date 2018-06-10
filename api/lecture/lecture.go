package lecture

import (
	"time"

	"github.com/kataras/iris"
	"github.com/mohuishou/scuplus-go/api"
	"github.com/mohuishou/scuplus-go/model"
)

func Get(ctx iris.Context) {
	page, err := ctx.URLParamInt("page")
	if err != nil {
		api.Error(ctx, 400, "参数错误", err)
		return
	}

	pageSize, err := ctx.URLParamInt("page_size")
	if err != nil {
		api.Error(ctx, 400, "参数错误", err)
		return
	}

	startTime := ctx.URLParam("start_time")
	scope := model.DB().Offset((page - 1) * pageSize).Limit(pageSize)

	if startTime == "" {
		scope = scope.Order("start_time desc")
	} else {
		t, err := time.Parse("2006-01-02", startTime)
		if err != nil {
			api.Error(ctx, 400, "参数错误", err)
			return
		}
		scope = scope.Where("start_time > ?", t).Order("start_time ASC")
	}

	var data []model.Lecture
	scope.Find(&data)
	api.Success(ctx, "获取成功", data)
}

type SearchParams struct {
	Page      int    `form:"page"`
	PageSize  int    `form:"page_size"`
	Title     string `form:"title"`
	StartTime string `form:"start_time"`
}

func Search(ctx iris.Context) {
	params := SearchParams{}
	ctx.ReadForm(&params)
	if params.Title == "" {
		api.Error(ctx, 400, "参数错误", nil)
		return
	}
	var data []model.Lecture
	scope := model.DB().Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize)
	if params.StartTime == "" {
		scope = scope.Order("start_time desc")
	} else {
		t, err := time.Parse("2006-01-02", params.StartTime)
		if err != nil {
			api.Error(ctx, 400, "参数错误", err)
			return
		}
		scope = scope.Where("start_time > ?", t).Order("start_time ASC")
	}
	if err := scope.Where("title like ?", "%"+params.Title+"%").Find(&data).Error; err != nil {
		api.Error(ctx, 90001, "获取错误", nil)
		return
	}
	api.Success(ctx, "获取成功！", data)
}
