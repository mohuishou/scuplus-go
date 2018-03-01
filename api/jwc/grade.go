package jwc

import "github.com/kataras/iris"
import "github.com/mohuishou/scuplus-go/middleware"
import "github.com/mohuishou/scuplus-go/model"

// GetGrade 获取成绩
func GetGrade(ctx iris.Context) {
	uid := middleware.GetUserID(ctx)

	ctx.JSON(map[string]interface{}{
		"status": 0,
		"data":   model.GetGrades(uid),
	})
}

// UpdateGrade 更新成绩
func UpdateGrade(ctx iris.Context) {
	uid := middleware.GetUserID(ctx)

	if _, err := model.UpdateGrades(uid); err != nil {
		ctx.JSON(map[string]interface{}{
			"status": 20001,
			"msg":    "更新失败！",
		})
		return
	}
	ctx.JSON(map[string]interface{}{
		"status": 0,
		"msg":    "更新成功",
		"data":   model.GetGrades(uid),
	})
}
