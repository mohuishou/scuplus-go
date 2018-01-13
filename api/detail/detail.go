package detail

import (
	"log"

	"github.com/kataras/iris"
	"github.com/mohuishou/scuplus-go/api"
	"github.com/mohuishou/scuplus-go/model"
)

// Params 参数
type Params struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Category string `form:"category"`
	TagID    uint   `form:"tag_id"`
}

// GetDetails 获取文章列表页
func GetDetails(ctx iris.Context) {
	params := Params{}
	ctx.ReadForm(&params)
	details := []model.Detail{}

	scope := model.DB().Select([]string{"details.id", "title", "url", "category", "details.created_at"})
	if params.Category != "" {
		scope = scope.Where("category = ?", params.Category)
	}

	scope = scope.Order("created_at desc").Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize)

	if params.TagID != 0 {
		// 获取tag
		tag := model.Tag{}
		model.DB().Find(&tag, params.TagID)
		scope = scope.Model(&tag).Related(&details, "Details")
		for _, v := range details {
			v.Tags = []model.Tag{tag}
		}
	} else {
		scope = scope.Preload("Tags").Find(&details)
	}
	err := scope.Error

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
			"page": params.Page,
			"data": details,
		},
	})
}

// GetDetail 获取文章详情
func GetDetail(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil || id == 0 {
		log.Println("[Error]: id:", id, err)
		api.Error(ctx, 40400, "参数错误", nil)
		return
	}
	detail := model.Detail{}
	model.DB().Preload("Tags").Find(&detail, id)
	api.Success(ctx, "获取成功！", detail)
}
