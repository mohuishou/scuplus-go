package route

import (
	"github.com/kataras/iris"
	"github.com/mohuishou/scuplus-go/api"
	"github.com/mohuishou/scuplus-go/api/course"
	"github.com/mohuishou/scuplus-go/api/detail"
	"github.com/mohuishou/scuplus-go/api/ecard"
	"github.com/mohuishou/scuplus-go/api/jwc"
	"github.com/mohuishou/scuplus-go/api/library"
	"github.com/mohuishou/scuplus-go/api/user"
	"github.com/mohuishou/scuplus-go/api/wechat"
	"github.com/mohuishou/scuplus-go/model"
)

// Routes 路由
func Routes(app *iris.Application) {
	app.Get("/notices", api.GetNotices)
	app.Get("/notice/{id}", api.GetNotice)
	app.Get("/user/grade", jwc.GetGrade)
	app.Post("/user/grade", jwc.UpdateGrade)
	app.Get("/user/schedule", jwc.GetSchedules)
	app.Post("/user/schedule", jwc.UpdateSchedule)
	app.Get("/user/exam", jwc.GetExam)
	app.Post("/user/exam", jwc.UpdateExam)
	app.Post("/user/feedback", user.FeedBack)
	app.Get("/user/feedbacks", user.GetFeedBacks)
	app.Get("/user/feedback/{id}", user.GetFeedBack)
	app.Post("/user/feedback/comment/{id}", user.FeedBackComment)
	app.Get("/user/ecard", ecard.Get)
	app.Post("/user/ecard", ecard.Update)
	app.Post("/user/msg_id", user.MsgID)
	app.Post("/login", api.Login)
	app.Post("/bind", api.Bind)
	app.Post("/library/bind", api.BindLibrary)
	app.Get("/details", detail.GetDetails)
	app.Get("/detail/tags", detail.GetTags)
	app.Get("/detail/{id}", detail.GetDetail)
	app.Post("/classroom", api.GetClassroom)
	app.Post("/library/search", library.Search)
	app.Get("/library/books", library.GetBook)
	app.Post("/library/loan", library.Loan)
	app.Get("/wechat/qcode", wechat.GetQCode)
	app.Get("/term", api.GetTerm)
	app.Get("/term/events", api.GetTermEvents)
	app.Post("/webhook", api.WebHook)

	// 课程相关api
	CourseRoutes(app)
}

// CourseRoutes 课程相关api
func CourseRoutes(app *iris.Application) {
	// 课程相关api需要绑定账号才能使用
	courseApp := app.Party("/course", func(ctx iris.Context) {
		uid := uint(ctx.Values().Get("user_id").(float64))
		u := model.User{}
		model.DB().Where("id = ? ", uid).Select([]string{"verify"}).Find(&u)
		if u.Verify == 0 {
			api.Error(ctx, 401, "用户尚未绑定！", nil)
			ctx.StopExecution()
			return
		}
		ctx.Next()
	})
	courseApp.Get("/", course.Get)
	courseApp.Get("s", course.GetCourses)
	courseApp.Post("/search", course.Search)
	courseApp.Post("/comment", course.Comment)
	courseApp.Post("/comment/update", course.UpdateComment)
	courseApp.Get("/comment", course.GetComment)
}
