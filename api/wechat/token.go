package wechat

import (
	"github.com/kataras/iris"
	"github.com/mohuishou/scuplus-go/api"
	"github.com/mohuishou/scuplus-go/util/wechat"
)

func Token(ctx iris.Context) {
	token, err := wechat.GetAccessToken(false)
	if err != nil {
		api.Error(ctx, 400, "token获取失败", err)
		return
	}
	api.Success(ctx, "获取成功！", token)
}
