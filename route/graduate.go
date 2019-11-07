package route

import (
	"github.com/kataras/iris/v12"
	"github.com/mohuishou/scuplus-go/api/graduate"
)

func GraduateRoutes(app *iris.Application) {
	graduateApp := app.Party("/graduate", VerifyBind)
	graduateApp.Get("/grades", graduate.GetGrade)
	graduateApp.Post("/grades", graduate.UpdateGrade)
	graduateApp.Get("/schedule", graduate.GetSchedule)
	graduateApp.Post("/schedule", graduate.UpdateSchedule)
}
