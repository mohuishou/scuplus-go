package model

import (
	"github.com/mohuishou/scu/jwc"
	"github.com/mohuishou/scu/jwc/evaluate"
	"github.com/mohuishou/scuplus-go/util"
)

const (
	// TermAutumn 秋季学期
	TermAutumn = 0
	// TermSpring 春季学期
	TermSpring = 1
)

type Evaluate struct {
	Model
	UserID       uint    `json:"user_id"`
	EvaluateID   string  `json:"evaluate_id"`
	EvaluateType string  `json:"evaluate_type"`
	TeacherName  string  `json:"teacher_name"`
	TeacherID    string  `json:"teacher_id"`
	TeacherType  int     `json:"teacher_type"` // 0: 老师, 1:助教
	CourseName   string  `json:"course_name"`
	CourseID     string  `json:"course_id"`
	Status       int     `json:"status"`
	Star         float64 `json:"star"`      //平均分
	Comment      string  `json:"comment"`   //评价内容
	LessonID     string  `json:"lesson_id"` // 课序号
	Year         int     `json:"year"`      // 年份
	Term         int     `json:"term"`      // 学期，0: 秋季学期, 1: 春季学期
}

// AfterCreate 新建之后回调
func (e *Evaluate) AfterCreate() error {
	if e.LessonID != "" {
		CBNewCourseEvaluate(e.UserID, e.CourseID, e.LessonID, e.CourseName)
	}
	return nil
}

// ConvertEvaluate 类型转换
func ConvertEvaluate(eva *evaluate.Evaluate) *Evaluate {
	return &Evaluate{
		EvaluateID:   eva.EvaluateID,
		EvaluateType: eva.EvaluateType,
		TeacherName:  eva.TeacherName,
		TeacherID:    eva.TeacherID,
		TeacherType:  eva.TeacherType,
		CourseName:   eva.CourseName,
		CourseID:     eva.CourseID,
		Status:       eva.Status,
		Star:         eva.Star,
		Comment:      eva.Comment,
	}
}

// UpdateEvaluateList 更新评教列表
func UpdateEvaluateList(uid uint) error {
	c, err := GetJwc(uid)
	if err != nil {
		return err
	}
	defer jwc.Logout(c)

	// 从教务处获取评教列表
	list, err := evaluate.GetEvaList(c)
	if err != nil {
		return err
	}

	for _, e := range list {
		eva := ConvertEvaluate(&e)
		eva.UserID = uid
		year, term := util.GetYearTerm()
		eva.Year = year
		eva.Term = term

		// 获取课序号
		cc := CourseCount{}
		DB().Where("course_id = ?", eva.CourseID).Where("teacher like ?", "%"+eva.TeacherName+"%").First(&cc)
		if cc.LessonID != "" {
			eva.LessonID = cc.LessonID
		}

		// 初始化，是否已存在，已存在则更新，不存在就新建
		old := Evaluate{}
		DB().FirstOrInit(&old, Evaluate{
			UserID:    uid,
			Year:      year,
			Term:      term,
			CourseID:  eva.CourseID,
			TeacherID: eva.TeacherID,
		})
		if old.ID == 0 {
			if err := DB().Create(eva).Error; err != nil {
				return err
			}
		} else {
			if err := DB().Model(&old).Updates(*eva).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

// DoEvaluate 评教
func DoEvaluate(eva *Evaluate, star float64, comment string) error {
	e := evaluate.Evaluate{
		EvaluateID:   eva.EvaluateID,
		EvaluateType: eva.EvaluateType,
		TeacherName:  eva.TeacherName,
		TeacherID:    eva.TeacherID,
		TeacherType:  eva.TeacherType,
		CourseName:   eva.TeacherName,
		CourseID:     eva.CourseID,
		Status:       eva.Status,
		Star:         star,
		Comment:      comment,
	}
	c, err := GetJwc(eva.UserID)
	if err != nil {
		return err
	}
	return evaluate.AddEvaluate(c, &e)
}
