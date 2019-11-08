package route

import (
	"github.com/kataras/iris/v12"
	"github.com/mohuishou/scuplus-go/api"
	"github.com/mohuishou/scuplus-go/api/library"
)

func LibraryRoutes(app *iris.Application) {
	app.Post("/library/bind", api.BindLibrary)

	libraryApp := app.Party("/library", VerifyLibrary)
	libraryApp.Post("/search", library.Search)
	libraryApp.Get("/books", library.GetBook)
	libraryApp.Post("/loan", library.Loan)
}
