package course

import (
	"fmt"

	"strings"

	"github.com/kataras/iris/v12"
	"github.com/mohuishou/scuplus-go/api"
	cache "github.com/mohuishou/scuplus-go/cache/lists"
	"github.com/mohuishou/scuplus-go/middleware"
	"github.com/mohuishou/scuplus-go/model"
)

// MinGradeAll 最少需要多少条统计
const MinGradeAll = 10

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

	rkey := fmt.Sprintf("courses.c%s.e%s.t%s.d%s.ca%s.o%s.p%d.ps%d", params.CallName, params.ExamType, params.Task, params.Day, params.Campus, params.Order, params.Page, params.PageSize)
	// 获取缓存信息
	data, err := cache.Get(rkey)
	if err == nil {
		ctx.Write(data)
		return
	}

	var courseCounts []model.CourseCount
	scope := model.DB().Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize)
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

	if params.Order == "" {
		params.Order = "avg_grade desc"
	}
	scope = scope.Order(params.Order).Order("good desc")
	if strings.Contains(params.Order, "avg_grade") {
		scope = scope.Where("grade_all > ?", MinGradeAll)
	}
	if strings.Contains(params.Order, "star") {
		scope = scope.Where("star > 0")
	}

	if err := scope.Find(&courseCounts).Error; err != nil {
		api.Error(ctx, 70001, "获取错误", nil)
		return
	}
	api.Success(ctx, "获取成功！", courseCounts)
	// 缓存数据,缓存12小时
	cache.Set(rkey, map[string]interface{}{
		"status": 0,
		"msg":    "获取成功！",
		"data":   courseCounts,
	}, 3600*12)
}

// SearchParams Get 参数
type SearchParams struct {
	Name        string `form:"name"`         // 搜索的课程名
	TeacherName string `form:"teacher_name"` // 搜索的教师名
	Order       string `form:"order"`        // 排序方式
	Page        int    `form:"page"`
	PageSize    int    `form:"page_size"`
}

// Search 课程搜索
func Search(ctx iris.Context) {
	params := SearchParams{}
	ctx.ReadForm(&params)
	if params.Name == "" && params.TeacherName == "" {
		api.Error(ctx, 70400, "参数错误", nil)
		return
	}
	var courseCounts []model.CourseCount
	scope := model.DB().Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize)

	if params.Name != "" {
		scope = scope.Where("name like ?", "%"+params.Name+"%")
	}

	if params.TeacherName != "" {
		scope = scope.Where("teacher like ?", "%"+params.TeacherName+"%")
	}

	if params.Order == "" {
		params.Order = "avg_grade desc"
	}
	scope = scope.Order(params.Order).Order("good desc")

	if err := scope.Find(&courseCounts).Error; err != nil {
		api.Error(ctx, 70001, "获取错误", nil)
		return
	}
	api.Success(ctx, "获取成功！", courseCounts)
}

// Get 获取一门课程的所有信息
// 包括CourseCount\Cousre\CousreEva中的所有信息
func Get(ctx iris.Context) {
	courseID := ctx.URLParam("course_id")
	lessonID := ctx.URLParam("lesson_id")
	if courseID == "" || lessonID == "" {
		api.Error(ctx, 70400, "参数错误", nil)
		return
	}

	// 获取课程信息
	var (
		courseCount     model.CourseCount
		courses         []model.Course
		courseEvaluates []model.CourseEvaluate
		courseGrades    []model.CourseGrade
	)
	scope := model.DB().Where("course_id = ? and lesson_id = ?", courseID, lessonID)
	scope.Find(&courseCount)
	scope.Find(&courses)
	scope.Find(&courseGrades)

	// todo: 获取用户昵称，用户头像,用户是否已经点赞
	scope.Where("status = 1").Order("updated_at desc").Find(&courseEvaluates)

	// 获取用户是否有该门课程
	uid := middleware.GetUserID(ctx)
	has := !model.DB().Where("user_id = ? and course_id = ? and lesson_id = ?", uid, courseID, lessonID).Select([]string{"id"}).Find(&model.Schedule{}).RecordNotFound()

	// 获取用户是否已经评价
	evaluate := model.CourseEvaluate{}
	model.DB().Where("user_id = ? and course_id = ? and lesson_id = ?", uid, courseID, lessonID).Select([]string{"id"}).Find(&evaluate)

	// 返回成功信息
	api.Success(ctx, "获取成功！", map[string]interface{}{
		"course_count":     courseCount,
		"courses":          courses,
		"course_evaluates": courseEvaluates,
		"course_grades":    courseGrades,
		"has":              has,      // true: 课程表中拥有该门课程, false: 不拥有
		"evaluate":         evaluate, // 是否已经评价
	})
}
