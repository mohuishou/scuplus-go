package detail

import (
	"log"
	"time"

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
	TagName  string `form:"tag_name"`
}

// GetDetails 获取文章列表页
func GetDetails(ctx iris.Context) {
	params := Params{}
	ctx.ReadForm(&params)
	if params.TagName != "" {
		t := model.Tag{}
		model.DB().Find(&t, "name = ?", params.TagName)
		params.TagID = t.ID
	}

	details := []model.Detail{}

	scope := model.DB().Select([]string{"details.id", "title", "url", "category", "details.created_at"})
	if params.Category != "" {
		scope = scope.Where("category = ?", params.Category)
	}

	scope = scope.Order("created_at desc").Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize)

	if params.TagName == "就业网" || params.TagName == "宣讲会" {
		// 获取tag
		tag := model.Tag{}
		model.DB().Find(&tag, params.TagID)
		scope = scope.Where("details.created_at > ?", time.Now().Format("2006-01-02")).Order("created_at asc").Model(&tag).Preload("Tags").Related(&details, "Details")
	} else if params.TagID != 0 {
		// 获取tag
		tag := model.Tag{}
		model.DB().Find(&tag, params.TagID)
		scope = scope.Model(&tag).Preload("Tags").Related(&details, "Details")
	} else if params.TagName == "全部" || (params.TagName == "" && params.TagID == 0) {
		scope = scope.Not("category", "就业网").Preload("Tags").Find(&details)
	} else {
		scope = scope.Preload("Tags").Find(&details)
	}
	err := scope.Error

	if err != nil {
		api.Error(ctx, 50001, "获取信息错误", nil)
		return
	}
	api.Success(ctx, "文章列表获取成功！", map[string]interface{}{
		"page": params.Page,
		"data": details,
	})
}

// GetDetail 获取文章详情
func GetDetail(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil || id == 0 {
		log.Println("[Error]: id:", id, err)
		api.Error(ctx, 50400, "参数错误", nil)
		return
	}
	detail := model.Detail{}
	model.DB().Preload("Tags").Find(&detail, id)
	api.Success(ctx, "获取成功！", detail)
}

// GetTags 获取所有标签
func GetTags(ctx iris.Context) {
	tags := []model.Tag{}
	model.DB().Find(&tags)
	api.Success(ctx, "获取成功！", tags)
}
