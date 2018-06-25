package graduate

import (
	"github.com/kataras/iris"
	"github.com/mohuishou/scuplus-go/api"
	"github.com/mohuishou/scuplus-go/middleware"
	"github.com/mohuishou/scuplus-go/model"
)

func GetGrade(ctx iris.Context) {
	uid := middleware.GetUserID(ctx)
	grades := []model.GraduateGrade{}
	model.DB().Where("user_id = ?", uid).Find(&grades)
	api.Success(ctx, "获取成功", grades)
}

func UpdateGrade(ctx iris.Context) {
	uid := middleware.GetUserID(ctx)
	if _, err := model.UpdateGraduateGrade(uid); err != nil {
		api.Error(ctx, 500, err.Error(), nil)
		return
	}
	api.Success(ctx, "更新成功", nil)
}
