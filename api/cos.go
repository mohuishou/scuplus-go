package api

import (
	"github.com/kataras/iris/v12"
	"github.com/mohuishou/scuplus-go/util/cos"
)

func COS(ctx iris.Context) {
	Success(ctx, "获取成功", cos.Sign())
}
