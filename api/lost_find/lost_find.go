package lostFind

import (
	"fmt"
	"log"

	"github.com/mohuishou/scuplus-go/model"

	"github.com/kataras/iris"
	"github.com/mohuishou/scuplus-go/api"
	cache "github.com/mohuishou/scuplus-go/cache/lists"
)

// ListParam 列表参数
type ListParam struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Category string `form:"category"`
}

// Lists 获取信息列表
func Lists(ctx iris.Context) {
	params := ListParams{}
	if err := ctx.ReadForm(&params); err != nil {
		api.Error(ctx, 80400, "参数错误", nil)
		return
	}

	// 获取缓存信息
	rkey := fmt.Sprintf("details.c%s.ps%d.p%d", params.Category, params.PageSize, params.Page)
	data, err := cache.Get(rkey)
	if err == nil {
		ctx.Write(data)
		return
	}

	var lists []model.LostFind
	// 获取列表信息
	model.DB().Offset((params.Page-1)*params.PageSize).Limit(params.PageSize).Where("category = ?", params.Category).Order("id desc").Find(&lists)
	api.Success(ctx, "获取成功！", lists)
	// 缓存数据,缓存半小时
	cache.Set(rkey, map[string]interface{}{
		"status": 0,
		"msg":    "获取成功！",
		"data":   lists,
	}, 3600*0.5)
}

// Get 获取一条信息
func Get(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil || id == 0 {
		log.Println("[Error]: id:", id, err)
		api.Error(ctx, 80400, "参数错误", nil)
		return
	}
	var data model.LostFind
	model.DB().Find(&data, id)
	api.Success(ctx,"获取成功！",map[string]interface{}{
		"data": data,
		"is_me": data.UserID == middleware.GetUserID(ctx),
	})
}

// NewParam 参数
type NewParam struct {
	ID uint `form:"id"`
	Title    string `form:"title" validate:"required"`    // 标题
	Pictures string `form:"pictures" validate:"required"` // 截图链接
	Info     string `form:"info" validate:"required,max=200"`     // 信息
	Address  string `form:"address" validate:"required,max=200"`  // 地点
	Contact  string `form:"contact" validate:"required,max=200"`  // 联系方式
	Category string `form:"category" validate:"required"` // 分类: 一卡通,其他,遗失
}

// Create 新建
func Create(ctx iris.Context) {
	data := param(ctx)
	if data == nil {
		return
	}

	// 不是一卡通直接保存
	if data.Category != "一卡通" {
		if err:=model.DB().Create(data).Error;err!=nil{
			api.Error(ctx,80001,"创建失败！",err)
			return
		}
		api.Success(ctx,"创建成功!",nil)
		return
	}

	// 一卡通，调用腾讯优图识别关键信息

}

// Update 更新一条信息
func Update(ctx iris.Context) {
	data := param(ctx)
	if data == nil  {
		return
	}
	if data.ID == 0 {
		api.Error(ctx, 80400, "参数错误！", err.Error())
		return nil 
	}

	var lost model.LostFind
	if err:=model.DB().Find(&lost,data.ID).Error;err!=nil{
		api.Error(ctx,80002,"更新失败",err)
		return
	}

	if err:=model.DB().Model(&lost).Updates(data).Error;err!=nil{
		api.Error(ctx,80003,"更新失败",err)
		return
	}
	api.Success(ctx,"更新成功！",nil)
}

func param(ctx iris.Context) *model.LostFind {
	params := NewParam{}
	if err := ctx.ReadForm(&params); err != nil{
		api.Error(ctx, 80400, "参数错误！", err)
		return nil
	}

	validate := validator.New()
	if err := validate.Struct(params); err != nil {
		api.Error(ctx, 80400, "参数错误！", err.Error())
		return nil 
	}

	return &model.LostFind{
		ID: params.ID,
		Title:params.Title,
		Category:params.Category,
		Info:params.Info,
		Address: params.Address,
		Contact: params.Contact,
		Pictures:params.Pictures,
	}
}
