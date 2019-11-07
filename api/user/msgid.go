package user

import (
	"github.com/kataras/iris/v12"
	"github.com/mohuishou/scuplus-go/api"
	"github.com/mohuishou/scuplus-go/cache/msgid"
	"github.com/mohuishou/scuplus-go/middleware"
)

func MsgID(context iris.Context) {
	uid := middleware.GetUserID(context)
	params := context.FormValues()
	ids, ok := params["ids"]
	if !ok {
		api.Error(context, 50400, "参数错误", nil)
		return
	}
	err := msgid.Set(uid, ids)
	if err != nil {
		api.Error(context, 50500, "缓存失败！", nil)
		return
	}
	api.Success(context, "缓存成功", nil)
}
