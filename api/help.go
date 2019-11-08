package api

import (
	"github.com/kataras/iris/v12"
	"github.com/mohuishou/scuplus-go/model"
)

// GetHelps 获取帮助信息
func GetHelps(ctx iris.Context) {
	items := []model.HelpItem{}
	model.DB().Where("sort > 0").Order("sort desc").Find(&items)
	Success(ctx, "获取成功！", items)
}
