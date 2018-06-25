package route

import (
	"github.com/kataras/iris"
	"github.com/mohuishou/scuplus-go/api"
	"github.com/mohuishou/scuplus-go/model"
)

// VerifyJWC 教务处绑定验证
func VerifyJWC(ctx iris.Context) {
	uid := GetUserID(ctx)
	u := model.User{}
	model.DB().Where("id = ? ", uid).Select([]string{"jwc_verify"}).Find(&u)
	if u.JwcVerify == 0 {
		api.Error(ctx, 401, "权限不足，请绑定本科教务处账号！", nil)
		ctx.StopExecution()
		return
	}
	ctx.Next()
}

// VerifyLibrary 图书馆绑定验证
func VerifyLibrary(ctx iris.Context) {
	uid := GetUserID(ctx)
	u := model.UserLibrary{}
	model.DB().Where("user_id = ? ", uid).Select([]string{"verify"}).Find(&u)
	if u.Verify == 0 {
		api.Error(ctx, 401, "权限不足，请绑定图书馆账号！", nil)
		ctx.StopExecution()
		return
	}
	ctx.Next()
}

// VerifyMy 个人中心绑定验证
func VerifyMy(ctx iris.Context) {
	uid := GetUserID(ctx)
	u := model.User{}
	model.DB().Where("id = ? ", uid).Select([]string{"verify"}).Find(&u)
	if u.Verify == 0 {
		api.Error(ctx, 401, "权限不足，请绑定信息门户账号！", nil)
		ctx.StopExecution()
		return
	}
	ctx.Next()
}

// VerifyBind 教务处，个人中心绑定任意一个账号都可以通过验证
func VerifyBind(ctx iris.Context) {
	uid := GetUserID(ctx)
	u := model.User{}
	model.DB().Where("id = ? ", uid).Select([]string{"jwc_verify", "verify"}).Find(&u)
	if u.JwcVerify == 0 && u.Verify == 0 {
		api.Error(ctx, 401, "权限不足，请绑定账号！", nil)
		ctx.StopExecution()
		return
	}
	ctx.Next()
}

// GetUserID 获取用户的id
func GetUserID(ctx iris.Context) uint {
	return uint(ctx.Values().Get("user_id").(float64))
}
