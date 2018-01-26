package main

import (
	"github.com/mohuishou/scuplus-go/config"
	"github.com/mohuishou/scuplus-go/middleware"

	"github.com/betacraft/yaag/irisyaag"
	"github.com/betacraft/yaag/yaag"
	"github.com/kataras/iris"
	"github.com/mohuishou/scuplus-go/route"
)

func main() {
	app := iris.New()

	yaag.Init(&yaag.Config{ // <- IMPORTANT, init the middleware.
		On:       true,
		DocTitle: "SCUPLUS",
		DocPath:  "apidoc/apidoc.html",
		BaseUrls: map[string]string{"Production": "", "Staging": ""},
	})
	app.Use(irisyaag.New()) // <- IMPORTANT, register the middleware.

	// 注册中间件
	middleware.Register(app)
	route.Routes(app)
	app.Run(iris.Addr("0.0.0.0:" + config.Get().Port))
}
