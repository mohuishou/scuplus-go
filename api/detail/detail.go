package detail

import (
	"github.com/kataras/iris"
	"github.com/mohuishou/scuplus-go/model"
)

// Params 参数
type Params struct {
	Page     int
	PageSize int
}

// GetDetails 获取文章
func GetDetails(ctx iris.Context) {
	params := Params{}
	ctx.ReadForm(&params)
	details := []model.Detail{}
	err := model.DB().Offset(params.Page * params.PageSize).Limit(params.PageSize).Find(&details).Error
	total := model.DB().Count(&model.Detail{})
	if err != nil {
		ctx.JSON(map[string]interface{}{
			"status":  20001,
			"message": "获取信息错误",
		})
		return
	}
	ctx.JSON(map[string]interface{}{
		"status":  0,
		"message": "获取成功",
		"data": map[string]interface{}{
			"page":  params.Page,
			"total": total,
			"data":  details,
		},
	})
}
