package route

import (
	"github.com/kataras/iris/v12"
	"github.com/mohuishou/scuplus-go/api/jwc"
)

// EvaluateRoutes 评教API
func EvaluateRoutes(app *iris.Application) {
	evaluateApp := app.Party("/", VerifyJWC)
	evaluateApp.Get("/evaluates", jwc.EvaluateList)
	evaluateApp.Post("/evaluates", jwc.UpdateEvaluateList)
	evaluateApp.Get("/evaluate/{id}", jwc.Evaluate)
	evaluateApp.Post("/evaluate", jwc.AddEvaluate)
}
