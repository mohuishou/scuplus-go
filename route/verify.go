package route

import (
	"github.com/kataras/iris/v12"
	"github.com/mohuishou/scuplus-go/api"
	"github.com/mohuishou/scuplus-go/cache/verify"
	"github.com/mohuishou/scuplus-go/model"
)

// VerifyJWC 教务处绑定验证
func VerifyJWC(ctx iris.Context) {
	if !getJwcVerify(ctx) {
		api.Error(ctx, 401, "权限不足，请绑定本科教务处账号！", nil)
		ctx.StopExecution()
		return
	}
	ctx.Next()
}

// VerifyLibrary 图书馆绑定验证
func VerifyLibrary(ctx iris.Context) {
	if !getLibraryVerify(ctx) {
		api.Error(ctx, 401, "权限不足，请绑定图书馆账号！", nil)
		ctx.StopExecution()
		return
	}
	ctx.Next()
}

// VerifyMy 个人中心绑定验证
func VerifyMy(ctx iris.Context) {
	if !getMyVerify(ctx) {
		api.Error(ctx, 401, "权限不足，请绑定信息门户账号！", nil)
		ctx.StopExecution()
		return
	}
	ctx.Next()
}

// VerifyBind 教务处，个人中心绑定任意一个账号都可以通过验证
func VerifyBind(ctx iris.Context) {
	if !getJwcVerify(ctx) && !getMyVerify(ctx) {
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

// 获取教务处和统一认证中心的验证
func getJwcVerify(ctx iris.Context) bool {
	uid := GetUserID(ctx)
	jwcVerify, _ := verify.Get(uid, "jwc")
	if jwcVerify {
		return jwcVerify
	}

	u := model.User{}
	model.DB().Where("id = ? ", uid).Select([]string{"jwc_verify"}).Find(&u)
	verify.Set(uid, "jwc", u.JwcVerify)
	return u.JwcVerify == 1
}

func getMyVerify(ctx iris.Context) bool {
	uid := GetUserID(ctx)
	myerify, _ := verify.Get(uid, "jwc")
	if myerify {
		return myerify
	}
	u := model.User{}
	model.DB().Where("id = ? ", uid).Select([]string{"verify"}).Find(&u)
	verify.Set(uid, "my", u.Verify)
	return u.Verify == 1
}

func getLibraryVerify(ctx iris.Context) bool {
	uid := GetUserID(ctx)
	v, _ := verify.Get(uid, "library")
	if v {
		return v
	}
	u := model.UserLibrary{}
	model.DB().Where("user_id = ? ", uid).Select([]string{"verify"}).Find(&u)
	verify.Set(uid, "library", u.Verify)
	return u.Verify == 1
}
