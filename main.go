package main

import (
	"github.com/mohuishou/scuplus-go/config"
	"github.com/mohuishou/scuplus-go/middleware"

	"github.com/mohuishou/scuplus-go/route"

	"github.com/kataras/iris"
)

func main() {
	app := iris.New()
	// 注册中间件
	middleware.Register(app)
	route.Routes(app)
	app.Run(iris.Addr("0.0.0.0:" + config.Get().Port))
}
