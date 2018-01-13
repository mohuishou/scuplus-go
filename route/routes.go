package route

import (
	"github.com/kataras/iris"
	"github.com/mohuishou/scuplus-go/api"
	"github.com/mohuishou/scuplus-go/api/detail"
	"github.com/mohuishou/scuplus-go/api/jwc"
)

// Routes 路由
func Routes(app *iris.Application) {
	app.Get("/user/grade", jwc.GetGrade)
	app.Post("/user/grade", jwc.UpdateGrade)
	app.Get("/user/shcedule", jwc.GetSchedules)
	app.Post("/user/shcedule", jwc.UpdateSchedule)
	app.Post("/login", api.Login)
	app.Post("/bind-jwc", api.BindJwc)
	app.Get("/details", detail.GetDetails)
	app.Get("/detail/{id}", detail.GetDetail)
}
