package route

import (
	"github.com/kataras/iris"
	"github.com/mohuishou/scuplus-go/api"
	"github.com/mohuishou/scuplus-go/api/detail"
	"github.com/mohuishou/scuplus-go/api/wechat"
)

// Routes 路由
func Routes(app *iris.Application) {
	app.Get("/notices", api.GetNotices)
	app.Get("/notice/{id}", api.GetNotice)

	app.Post("/login", api.Login)
	app.Post("/jwc/bind", api.BindJwc)
	app.Post("/bind", api.Bind)

	app.Get("/details", detail.GetDetails)
	app.Get("/detail/tags", detail.GetTags)
	app.Get("/detail/{id}", detail.GetDetail)
	app.Post("/classroom", api.GetClassroom)
	app.Get("/wechat/qcode", wechat.GetQCode)
	app.Get("/term", api.GetTerm)
	app.Get("/term/events", api.GetTermEvents)
	app.Get("/helps", api.GetHelps)
	UserRoutes(app)
	LibraryRoutes(app)
	// 课程相关api
	CourseRoutes(app)

	// 失物招领api
	LostFindRoutes(app)

	// 校园通讯录
	ContactRoutes(app)

	// 学术讲座
	LectureRoutes(app)

	// 评教
	EvaluateRoutes(app)

	// 研究生
	GraduateRoutes(app)
}
