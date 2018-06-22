package course

import (
	"log"

	jsoniter "github.com/json-iterator/go"
	"github.com/kataras/iris"
	"github.com/mohuishou/scuplus-go/api"
	"github.com/mohuishou/scuplus-go/middleware"
	"github.com/mohuishou/scuplus-go/model"
	"github.com/mohuishou/scuplus-go/util/wechat"
	validator "gopkg.in/go-playground/validator.v9"
)

// CommentParam 课程评价参数
type CommentParam struct {
	ID       uint   `form:"id"`
	CallName int    `form:"call_name" validate:"required,min=1,max=4"` // 点名方式
	ExamType int    `form:"exam_type" validate:"required,min=1,max=4"` // 考核方式
	Task     int    `form:"task" validate:"required,min=1,max=2"`      // 有无作业
	Star     int    `form:"star" validate:"required,min=1,max=3"`
	Comment  string `form:"comment" validate:"required,min=1,max=200"`
	//NickName string `form:"nick_name"`
	//Avatar   string `form:"avatar"`
}

// CommentListParams Get 参数
type CommentListParams struct {
	Name     string `form:"name"`
	Page     int    `form:"page" validate:"required,min=1"`
	PageSize int    `form:"page_size" validate:"required,min=1"`
}

// CommentList 获取评教列表
func CommentList(ctx iris.Context) {
	params := CommentListParams{}
	if err := ctx.ReadForm(&params); err != nil {
		api.Error(ctx, 70400, "参数获取错误！", err.Error())
		return
	}

	validate := validator.New()
	if err := validate.Struct(params); err != nil {
		api.Error(ctx, 70400, "参数校验错误！", err.Error())
		return
	}

	uid := middleware.GetUserID(ctx)

	courseEvaList := []model.CourseEvaluate{}
	scope := model.DB().Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize)
	if params.Name != "" {
		scope = scope.Where("course_name like ?", "%"+params.Name+"%")
	}
	scope.Where("user_id = ?", uid).Order("status asc").Find(&courseEvaList)
	api.Success(ctx, "获取成功！", courseEvaList)
}

// UpdateComment 更新评价
func UpdateComment(ctx iris.Context) {
	courseEvaluate := commentParam(ctx)
	if courseEvaluate == nil {
		return
	}

	if courseEvaluate.ID == 0 {
		api.Error(ctx, 70400, "参数错误！", nil)
		return
	}

	// 获取权限
	old := model.CourseEvaluate{}
	if err := model.DB().Find(&old, courseEvaluate.ID).Error; err != nil {
		api.Error(ctx, 70400, "参数错误！", nil)
		return
	}
	if old.UserID != courseEvaluate.UserID {
		api.Error(ctx, 70401, "参数错误！", nil)
		log.Println("用户权限错误！", courseEvaluate.UserID)
		return
	}

	if err := model.DB().Model(&old).Updates(courseEvaluate).Error; err != nil {
		api.Error(ctx, 70002, "更新失败！", err)
		return
	}
	api.Success(ctx, "更新成功！", nil)
}

func commentParam(ctx iris.Context) *model.CourseEvaluate {
	params := CommentParam{}
	if err := ctx.ReadForm(&params); err != nil {
		api.Error(ctx, 70400, "参数错误！", err)
		return nil
	}

	validate := validator.New()
	if err := validate.Struct(params); err != nil {
		api.Error(ctx, 70400, "参数错误！", err.Error())
		return nil
	}

	// 内容安全检查
	b, _ := jsoniter.Marshal(&params)
	res, err := wechat.MsgCheck(string(b))
	if !res {
		api.Error(ctx, 70005, "包含违法违规内容！", err)
		return nil
	}

	return &model.CourseEvaluate{
		Model:    model.Model{ID: params.ID},
		CallName: params.CallName,
		ExamType: params.ExamType,
		Task:     params.Task,
		Comment:  params.Comment,
		UserID:   middleware.GetUserID(ctx),
		Star:     params.Star,
		Status:   1,
	}
}

// GetComment 获取已经评价的课程
func GetComment(ctx iris.Context) {
	id, err := ctx.URLParamInt("id")
	if id == 0 || err != nil {
		api.Error(ctx, 70400, "参数错误！", err)
		return
	}
	courseEvaluate := model.CourseEvaluate{}
	if err := model.DB().Find(&courseEvaluate, id).Error; err != nil {
		api.Error(ctx, 70003, "获取失败！", err)
		return
	}

	if courseEvaluate.UserID != middleware.GetUserID(ctx) {
		api.Error(ctx, 70401, "您没有这个权限！", err)
		return
	}

	api.Success(ctx, "获取成功！", courseEvaluate)
}
