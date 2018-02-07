package jwc

import (
	"github.com/kataras/iris"
	"github.com/mohuishou/scuplus-go/api"
	"github.com/mohuishou/scuplus-go/middleware"
	"github.com/mohuishou/scuplus-go/model"
)

// GetExam 获取所有考表数据
func GetExam(ctx iris.Context) {
	uid := middleware.GetUserID(ctx)
	exams := []model.Exam{}
	model.DB().Where("user_id = ?", uid).Find(&exams)
	api.Success(ctx, "考表获取成功！", exams)
}

// UpdateExam 更新考表
func UpdateExam(ctx iris.Context) {
	uid := middleware.GetUserID(ctx)
	if err := model.UpdateExam(uid); err != nil {
		api.Error(ctx, 22001, "考表更新失败", nil)
	}
	api.Success(ctx, "更新成功!", nil)
}
