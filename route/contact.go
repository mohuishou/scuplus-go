package route

import (
	"github.com/kataras/iris/v12"
	"github.com/mohuishou/scuplus-go/api/contact"
)

// ContactRoutes 校园通讯录api
func ContactRoutes(app *iris.Application) {
	app.Get("/contact/categories", contact.Categories)
	app.Get("/contact/item/{cid}", contact.Get)
}
