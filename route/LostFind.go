package route

import (
	"github.com/kataras/iris/v12"
	"github.com/mohuishou/scuplus-go/api/lost_find"
)

// LostFindRoutes 失物招领api
func LostFindRoutes(app *iris.Application) {
	lostFindApp := app.Party("/", VerifyBind)
	lostFindApp.Post("/lost_find", lostFind.Create)
	lostFindApp.Get("/lost_finds", lostFind.Lists)
	lostFindApp.Get("/lost_find/{id}", lostFind.Get)
	lostFindApp.Post("/lost_find/update", lostFind.Update)
}
