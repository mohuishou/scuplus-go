package contact

import (
	"github.com/kataras/iris/v12"
	"github.com/mohuishou/scuplus-go/api"
	"github.com/mohuishou/scuplus-go/model"
)

// Categories 获取分类信息
func Categories(ctx iris.Context) {
	var data []model.ContactCategory
	model.DB().Find(&data)
	api.Success(ctx, "获取成功！", data)
}

// Get 获取某一分类下的所有通讯录
func Get(ctx iris.Context) {
	cid, err := ctx.Params().GetInt("cid")
	if err != nil || cid == 0 {
		api.Error(ctx, 90400, "参数错误！", err)
		return
	}
	var data = []model.ContactBook{}
	model.DB().Where("contact_category_id = ?", cid).Find(&data)
	api.Success(ctx, "获取成功！", data)
}
