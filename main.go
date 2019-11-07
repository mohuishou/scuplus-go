package main

import (
	"github.com/kataras/iris/v12/middleware/pprof"
	"github.com/mohuishou/scuplus-go/config"
	"github.com/mohuishou/scuplus-go/middleware"

	"github.com/kataras/iris/v12"
	"github.com/mohuishou/scuplus-go/route"
)

func main() {
	app := iris.New()

	// 注册中间件
	middleware.Register(app)
	route.Routes(app)
	app.Any("/debug/pprof/{action:path}", pprof.New())
	app.Run(iris.Addr("0.0.0.0:" + config.Get().Port))
}
