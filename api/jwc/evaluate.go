package jwc

import (
	"log"

	jsoniter "github.com/json-iterator/go"
	"github.com/kataras/iris"
	"github.com/mohuishou/scuplus-go/api"
	"github.com/mohuishou/scuplus-go/middleware"
	"github.com/mohuishou/scuplus-go/model"
	"github.com/mohuishou/scuplus-go/util"
	"github.com/mohuishou/scuplus-go/util/wechat"
	"gopkg.in/go-playground/validator.v9"
)

type ListParams struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

// EvaluateList list
func EvaluateList(ctx iris.Context) {
	uid := middleware.GetUserID(ctx)

	params := ListParams{}

	if err := ctx.ReadForm(&params); err != nil {
		api.Error(ctx, 84400, "参数错误！", err)
		return
	}
	// Todo: 前端有个bug，暂时将所有课程返回，新版本之后再去除掉
	//params.PageSize = 100

	msg := "获取成功"

	// 本学期评教列表是否存在，不存在则尝试获取教务处
	year, term := util.GetYearTerm()
	notFound := model.DB().Where("user_id = ? and year = ? and term = ?", uid, year, term).Find(&model.Evaluate{}).RecordNotFound()
	if notFound && model.UpdateEvaluateList(uid) != nil {
		msg = "评教列表获取失败"
	}

	// 获取评教列表
	var evaList []model.Evaluate
	scope := model.DB().Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize)
	scope = scope.Select([]string{"id", "teacher_name", "evaluate_type", "course_name", "status", "year", "term"})
	scope.Where("user_id = ?", uid).Order("year desc").Order("term desc").Find(&evaList)
	api.Success(ctx, msg, evaList)
}

// EvaluateList 手动更新评教列表
func UpdateEvaluateList(ctx iris.Context) {
	uid := middleware.GetUserID(ctx)
	if err := model.UpdateEvaluateList(uid); err != nil {
		api.Error(ctx, 84005, "更新失败", err)
		return
	}
	api.Success(ctx, "更新成功！", nil)
}

func Evaluate(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	uid := middleware.GetUserID(ctx)
	if err != nil || id == 0 {
		log.Println("[Error]: id:", id, err)
		api.Error(ctx, 84400, "参数错误", nil)
		return
	}

	// 获取教务处评教信息
	eva := model.Evaluate{}
	if err := model.DB().Find(&eva, id).Error; err != nil {
		api.Error(ctx, 84000, "数据获取失败", err)
		return
	}

	// 获取评教信息
	courseEva := model.CourseEvaluate{}
	if eva.LessonID != "" {
		model.DB().Where("user_id = ? and course_id = ? and lesson_id = ?", uid, eva.CourseID, eva.LessonID).Find(&courseEva)
	}
	api.Success(ctx, "获取成功", map[string]interface{}{
		"evaluate":        eva,
		"course_evaluate": courseEva,
	})

}

type AddParams struct {
	ID       uint    `form:"id" validate:"required"`
	CallName int     `form:"call_name" validate:"required,min=1,max=4"` // 点名方式
	ExamType int     `form:"exam_type" validate:"required,min=1,max=4"` // 考核方式
	Task     int     `form:"task" validate:"required,min=1,max=2"`      // 有无作业
	Star     float64 `form:"star" validate:"required,min=1,max=3"`
	Comment  string  `form:"comment" validate:"required,min=1,max=200"`
}

// AddEvaluate 添加评教信息
func AddEvaluate(ctx iris.Context) {
	// 参数校验
	params := AddParams{}
	if err := ctx.ReadForm(&params); err != nil {
		api.Error(ctx, 84400, "参数获取错误！", err)
		return
	}

	validate := validator.New()
	if err := validate.Struct(params); err != nil {
		api.Error(ctx, 84400, "参数校验错误！", err.Error())
		return
	}

	// 内容安全检查
	b, _ := jsoniter.Marshal(&params)
	res, err := wechat.MsgCheck(string(b))
	if !res {
		api.Error(ctx, 70005, "包含违法违规内容！", err)
		return
	}
	// 前端的评分为321,整体加两分
	params.Star = params.Star + 2

	eva := model.Evaluate{}
	if err := model.DB().Find(&eva, params.ID).Error; err != nil {
		api.Error(ctx, 70005, "数据不存在！", err)
		return
	}

	// 教务处评教
	if err := model.DoEvaluate(&eva, params.Star, params.Comment); err != nil {
		api.Error(ctx, 84400, err.Error(), nil)
		return
	}

	// 更新数据库
	model.DB().Model(&eva).Updates(map[string]interface{}{
		"star":    params.Star,
		"comment": params.Comment,
	})

	if eva.LessonID == "" {
		api.Success(ctx, "评教成功", nil)
		return
	}

	// we川大评教
	courseEva := model.CourseEvaluate{
		Star:     int(params.Star - 2),
		CallName: params.CallName,
		ExamType: params.ExamType,
		Comment:  params.Comment,
		Task:     params.Task,
		Status:   1,
		Score:    1,
	}
	oldCourseEva := model.CourseEvaluate{}
	model.DB().Where(
		"user_id = ? and course_id = ? and lesson_id = ?",
		eva.UserID,
		eva.CourseID,
		eva.LessonID,
	).Find(&oldCourseEva)
	model.DB().Model(&oldCourseEva).Updates(courseEva)
	api.Success(ctx, "评教成功", nil)
}
