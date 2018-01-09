package jwc

import (
	"github.com/kataras/iris"
	"github.com/mohuishou/scuplus-go/middleware"
	"github.com/mohuishou/scuplus-go/model"
)

// GetSchedules 获取课程表
func GetSchedules(ctx iris.Context) {
	uid := middleware.GetUserID(ctx)
	term := ctx.FormValue("term")

	ctx.JSON(map[string]interface{}{
		"status": 0,
		"data":   model.GetSchedules(uid, term),
	})
}

// UpdateSchedule 更新课程表
func UpdateSchedule(ctx iris.Context) {
	uid := middleware.GetUserID(ctx)
	term := ctx.FormValue("term")
	if err := model.UpdateSchedules(uid, term); err != nil {
		ctx.JSON(map[string]interface{}{
			"status":  10001,
			"message": "更新失败！",
		})
		return
	}
	ctx.JSON(map[string]interface{}{
		"status":  0,
		"message": "更新成功",
		"data":    model.GetSchedules(uid, term),
	})
}
