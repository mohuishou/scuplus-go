package route

import (
	"github.com/kataras/iris/v12"
	"github.com/mohuishou/scuplus-go/api/course"
)

// CourseRoutes 课程相关api
func CourseRoutes(app *iris.Application) {
	// 课程相关api需要绑定账号才能使用
	courseApp := app.Party("/course", VerifyJWC)
	courseApp.Get("/", course.Get)
	courseApp.Get("/all", course.GetCourses)
	courseApp.Post("/search", course.Search)
	courseApp.Post("/comment", course.UpdateComment)
	courseApp.Get("/comment", course.GetComment)
	courseApp.Get("/comments", course.CommentList)
}
