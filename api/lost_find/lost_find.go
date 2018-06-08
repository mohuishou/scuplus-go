package lostFind

import (
	"fmt"
	"log"

	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/mohuishou/scuplus-go/job"

	"github.com/mohuishou/scuplus-go/model"

	"github.com/json-iterator/go"
	"github.com/kataras/iris"
	"github.com/mohuishou/scuplus-go/api"
	cache "github.com/mohuishou/scuplus-go/cache/lists"
	"github.com/mohuishou/scuplus-go/middleware"
	"gopkg.in/go-playground/validator.v9"
)

// ListParam 列表参数
type ListParam struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	My       int    `form:"my"`
	Category string `form:"category"`
}

// Lists 获取信息列表
func Lists(ctx iris.Context) {
	params := ListParam{}
	if err := ctx.ReadForm(&params); err != nil {
		api.Error(ctx, 80400, "参数错误", err)
		return
	}

	// 获取缓存信息,个人列表不需要缓存
	rkey := fmt.Sprintf("details.c%s.ps%d.p%d", params.Category, params.PageSize, params.Page)
	if params.My == 0 {
		data, err := cache.Get(rkey)
		if err == nil {
			ctx.Write(data)
			return
		}
	}

	var lists []model.LostFind
	fields := []string{"id", "title", "created_at", "nickname"}
	if params.Category != model.LostFindCard {
		fields = append(fields, "pictures")
	}
	// 获取列表信息
	scope := model.DB().Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize)
	scope = scope.Order("id desc").Select(fields)
	scope = scope.Where("status = 1").Where("category = ?", params.Category)
	if params.My == 1 {
		scope = scope.Where("user_id = ?", middleware.GetUserID(ctx))
	}
	scope.Find(&lists)
	api.Success(ctx, "获取成功！", lists)

	if params.My == 0 {
		// 缓存数据,缓存半小时
		cache.Set(rkey, map[string]interface{}{
			"status": 0,
			"msg":    "获取成功！",
			"data":   lists,
		}, 3600*0.5)
	}

}

// Get 获取一条信息
func Get(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil || id == 0 {
		log.Println("[Error]: id:", id, err)
		api.Error(ctx, 80400, "参数错误", nil)
		return
	}

	// 获取数据
	var data model.LostFind
	if err := model.DB().Find(&data, id).Error; err != nil {
		api.Error(ctx, 80500, "数据获取失败！", err)
		return
	}

	isOwner := false
	uid := middleware.GetUserID(ctx)
	// 判断一卡通是否为拥有者
	if data.Category == model.LostFindCard {
		cardInfo := map[string]string{}
		jsoniter.Unmarshal([]byte(data.CardInfo), &cardInfo)
		if no, ok := cardInfo["no"]; ok {
			owner := model.User{}
			model.DB().Where("jwc_student_id = ?", no).Find(&owner)
			if owner.ID == uid {
				isOwner = true
			}
		}
	}

	api.Success(ctx, "获取成功！", map[string]interface{}{
		"data":     data,
		"is_me":    data.UserID == uid,
		"is_owner": isOwner,
	})
}

// NewParam 参数
type NewParam struct {
	ID       uint   `form:"id"`
	Title    string `form:"title" validate:"required"`           // 标题
	Pictures string `form:"pictures"`                            // 截图链接
	Info     string `form:"info" validate:"required,max=200"`    // 信息
	Address  string `form:"address" validate:"required,max=200"` // 地点
	Contact  string `form:"contact" validate:"required,max=200"` // 联系方式
	Category string `form:"category" validate:"required"`        // 分类: 一卡通,其他,遗失
	Nickname string `form:"nickname" validate:"required,max=200"`
}

// Create 新建
func Create(ctx iris.Context) {
	data := param(ctx)
	if data == nil {
		return
	}

	// 一卡通需要识别之后再做展示
	data.Status = model.LostFindShow
	if data.Category == model.LostFindCard {
		data.Status = model.LostFindReady
	}

	res := model.DB().Create(data)
	if err := res.Error; err != nil {
		api.Error(ctx, 80001, "创建失败！", err)
		return
	}

	if data.Category == model.LostFindCard {
		// 一卡通，异步调用腾讯优图识别关键信息
		sign := &tasks.Signature{
			Name: "card_ocr",
			Args: []tasks.Arg{
				{
					Type:  "uint",
					Value: data.ID,
				},
			},
		}
		_, err := job.Server.SendTask(sign)
		if err != nil {
			log.Println("ocr create err", err)
		}
	}

	api.Success(ctx, "创建成功!", nil)
	return
}

// Update 更新一条信息
func Update(ctx iris.Context) {
	data := param(ctx)
	if data == nil {
		return
	}
	if data.ID == 0 {
		api.Error(ctx, 80400, "参数错误！", nil)
		return
	}

	var lost model.LostFind
	if err := model.DB().Find(&lost, data.ID).Error; err != nil {
		api.Error(ctx, 80002, "更新失败", err)
		return
	}

	// 一卡通图片更改需要重新识别
	isOCR := false
	if data.Category == model.LostFindCard && data.Pictures != lost.Pictures {
		data.Status = 0
		data.CardInfo = ""
		isOCR = true
	}

	if err := model.DB().Save(&data).Error; err != nil {
		api.Error(ctx, 80003, "更新失败", err)
		return
	}

	if isOCR {
		// 一卡通，异步调用腾讯优图识别关键信息
		sign := &tasks.Signature{
			Name: "card_ocr",
			Args: []tasks.Arg{
				{
					Type:  "uint",
					Value: data.ID,
				},
			},
		}
		_, err := job.Server.SendTask(sign)
		if err != nil {
			log.Println("ocr update err", err)
		}
	}
	api.Success(ctx, "更新成功！", nil)
}

func param(ctx iris.Context) *model.LostFind {
	params := NewParam{}
	if err := ctx.ReadForm(&params); err != nil {
		api.Error(ctx, 80400, "参数错误！", err)
		return nil
	}

	validate := validator.New()
	if err := validate.Struct(params); err != nil {
		api.Error(ctx, 80400, "参数错误！", err.Error())
		return nil
	}

	return &model.LostFind{
		Model:    model.Model{ID: params.ID},
		Title:    params.Title,
		Category: params.Category,
		Nickname: params.Nickname,
		Info:     params.Info,
		Address:  params.Address,
		Contact:  params.Contact,
		Pictures: params.Pictures,
		UserID:   middleware.GetUserID(ctx),
	}
}
