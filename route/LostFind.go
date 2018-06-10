package route

import (
	"github.com/kataras/iris"
	"github.com/mohuishou/scuplus-go/api/lost_find"
)

// LostFindRoutes 失物招领api
func LostFindRoutes(app *iris.Application) {
	app.Post("/lost_find", lostFind.Create)
	app.Get("/lost_finds", lostFind.Lists)
	app.Get("/lost_find/{id}", lostFind.Get)
	app.Post("/lost_find/update", lostFind.Update)
}
