// iris provides some basic middleware, most for your learning courve.
// You can use any net/http compatible middleware with iris.FromStd wrapper.
//
// JWT net/http video tutorial for golang newcomers: https://www.youtube.com/watch?v=dgJFeqeXVKw
//
// This middleware is the only one cloned from external source: https://github.com/auth0/go-jwt-middleware
// (because it used "context" to define the user but we don't need that so a simple iris.FromStd wouldn't work as expected.)
package main

// $ go get -u github.com/dgrijalva/jwt-go
// $ go run main.go

import (
	"github.com/mohuishou/scuplus-go/middleware"

	"github.com/mohuishou/scuplus-go/route"

	"github.com/kataras/iris"

	"github.com/dgrijalva/jwt-go"
)

func myHandler(ctx iris.Context) {
	user := ctx.Values().Get("jwt").(*jwt.Token)

	ctx.Writef("This is an authenticated request\n")
	ctx.Writef("Claim content:\n")

	ctx.Writef("%s", int64(user.Claims.(jwt.MapClaims)["nbf"].(float64)))
}

func main() {
	app := iris.New()

	// 注册中间件
	middleware.Register(app)

	app.Get("/ping", myHandler)
	route.Routes(app)
	app.Run(iris.Addr("localhost:80"))
} // don't forget to look ../jwt_test.go to seee how to set your own custom claims
