package course

import (
	"github.com/kataras/iris"
	"github.com/mohuishou/scuplus-go/api"
	"github.com/mohuishou/scuplus-go/model"
)

// GetParams Get 参数
type GetParams struct {
	CallName string `form:"call_name"` // 点名方式
	ExamType string `form:"exam_type"` // 考核方式
	Task     string `form:"task"`      // 有无作业
	Day      string `form:"day"`       // 周几上课
	Campus   string `form:"campus"`    // 校区
	Order    string `form:"order"`     // 排序
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

// GetCourses 获取课程列表
func GetCourses(ctx iris.Context) {
	params := GetParams{}
	ctx.ReadForm(&params)
	var courseCounts []model.CourseCount
	scope := model.DB().Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize).Order("avg_grade desc")
	if params.CallName != "" {
		scope = scope.Where("call_name = ?", params.CallName)
	}
	if params.Day != "" {
		scope = scope.Where("day = ?", params.Day)
	}
	if params.ExamType != "" {
		scope = scope.Where("exam_type = ?", params.ExamType)
	}
	if params.Task != "" {
		scope = scope.Where("task = ?", params.Task)
	}
	if params.Campus != "" {
		scope = scope.Where("campus = ?", params.Campus)
	}
	if params.Order != "" {
		scope = scope.Order(params.Order)
	}
	if err := scope.Find(&courseCounts).Error; err != nil {
		api.Error(ctx, 70001, "获取错误", nil)
		return
	}
	api.Success(ctx, "获取成功！", courseCounts)
}

// SearchParams Get 参数
type SearchParams struct {
	Name     string `form:"name"` // 搜索的课程名
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

// Search 课程搜索
func Search(ctx iris.Context) {
	params := SearchParams{}
	ctx.ReadForm(&params)
	if params.Name == "" {
		api.Error(ctx, 70400, "参数错误", nil)
		return
	}
	var courseCounts []model.CourseCount
	scope := model.DB().Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize).Order("avg_grade desc")
	if err := scope.Where("%name% = ?", params.Name).Find(&courseCounts).Error; err != nil {
		api.Error(ctx, 70001, "获取错误", nil)
		return
	}
	api.Success(ctx, "获取成功！", courseCounts)
}

// Get 获取一门课程的所有信息
// 包括CourseCount\Cousre\CousreEva中的所有信息
func Get() {

}

// Comment 课程评价，目前只能评价正在上的课程
func Comment() {

}

// GetComment 获取已经评价的课程
func GetComment() {

}
