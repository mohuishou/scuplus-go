package route

import (
	"github.com/kataras/iris"
	"github.com/mohuishou/scuplus-go/api/jwc"
)

// EvaluateRoutes 评教API
func EvaluateRoutes(app *iris.Application) {
	app.Get("/evaluates", jwc.EvaluateList)
	app.Post("/evaluates", jwc.UpdateEvaluateList)
	app.Get("/evaluate/{id}", jwc.Evaluate)
	app.Post("/evaluate", jwc.AddEvaluate)
}
