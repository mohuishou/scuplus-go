package route

import (
	"github.com/kataras/iris/v12"
	"github.com/mohuishou/scuplus-go/api/ecard"
	"github.com/mohuishou/scuplus-go/api/jwc"
	"github.com/mohuishou/scuplus-go/api/user"
)

func UserRoutes(app *iris.Application) {
	// 需要教务处验证的api
	userVerifyJWCApp := app.Party("/user", VerifyJWC)
	userVerifyJWCRoutes(userVerifyJWCApp)

	// 需要信息门户认证的api
	userVerifyMyApp := app.Party("/user", VerifyMy)
	userVerifyMyRoutes(userVerifyMyApp)

	// 不需要验证的api
	userApp := app.Party("/user")
	userRoutes(userApp)
}

func userVerifyMyRoutes(app iris.Party) {
	app.Get("/ecard", ecard.Get)
	app.Post("/ecard", ecard.Update)
}

func userRoutes(app iris.Party) {
	app.Post("/feedback", user.FeedBack)
	app.Get("/feedbacks", user.GetFeedBacks)
	app.Get("/feedback/{id}", user.GetFeedBack)
	app.Post("/feedback/comment/{id}", user.FeedBackComment)
	app.Post("/config/notify", user.UpdateNotify)
	app.Get("/config/notify", user.GetNotify)
	app.Post("/config/type", user.UpdateUserType)
	app.Post("/msg_id", user.MsgID)
}

func userVerifyJWCRoutes(app iris.Party) {
	app.Get("/grade", jwc.GetGrade)
	app.Post("/grade", jwc.UpdateGrade)
	app.Get("/schedule", jwc.GetSchedules)
	app.Post("/schedule", jwc.UpdateSchedule)
	app.Get("/exam", jwc.GetExam)
	app.Post("/exam", jwc.UpdateExam)
}
