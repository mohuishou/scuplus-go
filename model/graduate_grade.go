package model

import (
	"errors"

	"log"

	"github.com/deckarep/golang-set"
	"github.com/jinzhu/gorm"
	"github.com/mohuishou/scu/ehall/grade"
)

// GraduateGrade 研究生成绩表
type GraduateGrade struct {
	Model
	UserID     uint    `json:"user_id"`
	CourseID   string  `json:"course_id"`   // 课程代码
	CourseName string  `json:"course_name"` // 课程名
	Credit     float64 `json:"credit"`      // 学分
	CourseType string  `json:"course_type"` // 课程类别
	GradeShow  string  `json:"grade_show"`  // 成绩显示值
	GPA        float64 `json:"gpa"`         // 绩点
	Grade      float64 `json:"grade"`       // 百分制成绩
	ExamType   string  `json:"exam_type"`   // 考试性质
	YearTerm   string  `json:"year_term"`   // 学年学期
	Term       int     `json:"term"`        //0: 秋季学期, 1: 春季学期
	Year       int     `json:"year"`
}

// UpdateGraduateGrade 更新研究生用户成绩
// 返回更新的数据slice, 如果为空则表示没有新的成绩
// 操作步骤: 抓取用户成绩(N) -> 获取数据用户已有成绩(M) -> 数据库删除数据(M-N) -> 数据库插入数据(N-M)
// 用户成绩更新大多数情况都是新增，对于只是更新个别字段的成绩直接删除即可
func UpdateGraduateGrade(uid uint) ([]GraduateGrade, error) {
	// 1. 登录用户统一认证中心
	c, err := GetCollector(uid)
	if err != nil {
		return nil, err
	}

	// 2. 获取用户成绩
	grades, err := grade.Get(c)
	if err != nil {
		return nil, err
	}
	if len(grades) == 0 {
		return nil, errors.New("没有获取到成绩信息！")
	}

	// 3. slice to set
	newGradeSet := mapset.NewSet()
	for _, g := range grades {
		tmpGrade := &GraduateGrade{}
		tmpGrade.convert(g)
		newGradeSet.Add(*tmpGrade)
	}

	// 4. 从数据库取出现有数据
	oldGradeSet, oldGradeIDs, err := oldGraduateGradeSet(uid)
	if err != nil {
		return nil, err
	}

	// 5. 新建事务，准备开始数据库操作
	tx := DB().Begin()

	// 6. 删除需要删除的数据
	deleteSet := oldGradeSet.Difference(newGradeSet)
	deleteIDs := graduateGradeDeleteIDs(deleteSet, oldGradeIDs)
	if len(deleteIDs) != 0 {
		tx.Unscoped().Where(deleteIDs).Delete(GraduateGrade{})
	}

	// 7. 新增更新的数据
	createSet := newGradeSet.Difference(oldGradeSet)
	it := createSet.Iterator()
	var updateGrades []GraduateGrade
	for elem := range it.C {
		g := elem.(GraduateGrade)
		g.UserID = uid
		updateGrades = append(updateGrades, g)
		if err := tx.Create(&g).Error; err != nil {
			tx.Rollback()
			log.Println("[Error]: UpdateGrades", err)
			return nil, err
		}
	}
	// 提交事务
	tx.Commit()
	return updateGrades, nil
}

// graduateGradeDeleteIDs 获取需要删除数据的id
func graduateGradeDeleteIDs(deleteSet mapset.Set, oldGradeIDs map[string]uint) []uint {
	it := deleteSet.Iterator()
	var deleteIDs []uint
	for elem := range it.C {
		g := elem.(GraduateGrade)
		if id, ok := oldGradeIDs[g.CourseID]; ok {
			deleteIDs = append(deleteIDs, id)
		}
	}
	return deleteIDs
}

// oldGraduateGradeSet 获取已有的成绩数据，并将其转换为set
func oldGraduateGradeSet(uid uint) (mapset.Set, map[string]uint, error) {
	var oldGrades []GraduateGrade
	err := DB().Where("user_id = ?", uid).Find(&oldGrades).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, nil, err
	}
	// slice to set
	oldGradeSet := mapset.NewSet()
	oldGradeIDs := map[string]uint{}
	for _, v := range oldGrades {
		oldGradeIDs[v.CourseID] = v.ID
		v.UserID = 0
		v.Model = Model{}
		oldGradeSet.Add(v)
	}
	return oldGradeSet, oldGradeIDs, nil
}

// convert 数据类型转换
func (g *GraduateGrade) convert(data grade.Grade) {
	g.CourseID = data.CourseID
	g.CourseName = data.CourseName
	g.Credit = data.Credit
	g.CourseType = data.CourseType
	g.GradeShow = data.GradeShow
	g.GPA = data.GPA
	g.Grade = data.Grade
	g.GradeShow = data.GradeShow
	g.ExamType = data.ExamType
	g.YearTerm = data.YearTerm
	g.Term = data.Term
	g.Year = data.Year
}
