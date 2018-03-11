package api

import (
	"github.com/kataras/iris"
	"github.com/mohuishou/scuplus-go/util/github"
)

func WebHook(ctx iris.Context) {
	github.WebHook(ctx.Request())
	ctx.JSON("ok")
}
