package route

import (
	"github.com/kataras/iris"
	"github.com/mohuishou/scuplus-go/api"
	"github.com/mohuishou/scuplus-go/api/course"
	"github.com/mohuishou/scuplus-go/model"
)

// CourseRoutes 课程相关api
func CourseRoutes(app *iris.Application) {
	// 课程相关api需要绑定账号才能使用
	courseApp := app.Party("/course", func(ctx iris.Context) {
		uid := uint(ctx.Values().Get("user_id").(float64))
		u := model.User{}
		model.DB().Where("id = ? ", uid).Select([]string{"jwc_verify"}).Find(&u)
		if u.JwcVerify == 0 {
			api.Error(ctx, 401, "用户尚未绑定！", nil)
			ctx.StopExecution()
			return
		}
		ctx.Next()
	})
	courseApp.Get("/", course.Get)
	courseApp.Get("/all", course.GetCourses)
	courseApp.Post("/search", course.Search)
	courseApp.Post("/comment", course.UpdateComment)
	courseApp.Get("/comment", course.GetComment)
	courseApp.Get("/comments", course.CommentList)
}
