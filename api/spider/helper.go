package spider

import (
	"github.com/kataras/iris"
	"github.com/mohuishou/scuplus-go/util/spider/helper/jwc"
)

// GetJwcCookies 获取教务处cookies
func GetJwcCookies(ctx iris.Context) {
	cookies := jwc.GetCookies()
	ctx.JSON(cookies)
}
