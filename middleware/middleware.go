package middleware

import (
	"github.com/kataras/iris/v12"
)

// Register 注册中间件
func Register(app *iris.Application) {
	app.Use(jwtMiddle)
}
