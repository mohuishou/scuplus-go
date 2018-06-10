package route

import (
	"github.com/kataras/iris"
	"github.com/mohuishou/scuplus-go/api/lecture"
)

func LectureRoutes(app *iris.Application) {
	app.Get("/lectures", lecture.Get)
	app.Post("/lectures/search", lecture.Search)
}
