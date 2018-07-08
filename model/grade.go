package model

import (
	"errors"
	"log"

	"fmt"

	"github.com/deckarep/golang-set"
	"github.com/mohuishou/scu/jwc/grade"
)

// Grade 成绩
type Grade struct {
	Model
	UserID     uint   `json:"user_id" gorm:"unique_index:u_grade"`
	CourseID   string `json:"course_id" gorm:"unique_index:u_grade"`
	LessonID   string `json:"lesson_id" gorm:"unique_index:u_grade"`
	CourseName string `json:"course_name"`
	Credit     string `json:"credit"`
	CourseType string `json:"course_type"`
	Grade      string `json:"grade"`
	Term       int    `json:"term" gorm:"unique_index:u_grade"` //0: 秋季学期, 1: 春季学期
	Year       int    `json:"year" gorm:"unique_index:u_grade"`
	TermName   string `json:"term_name"`
}

// AfterCreate 新建之后回调
func (g *Grade) AfterCreate() error {
	CBNewCourseEvaluate(g.UserID, g.CourseID, g.LessonID, g.CourseName)
	return nil
}

// GetGrades 获取用户的所有成绩
func GetGrades(userID uint) []Grade {
	grades := make([]Grade, 0)
	failGrades := make([]Grade, 0)
	DB().Where("user_id = ? and year = -1", userID).Order("year desc").Find(&failGrades)
	if err := DB().Where("user_id = ? and year != -1", userID).Order("year desc, term desc").Find(&grades).Error; err != nil {
		log.Printf("[Error] GetGrades Fail, userID: %d, err: %s", userID, err.Error())
	}
	return append(failGrades, grades...)
}

// UpdateGrades 更新用户的所有成绩
// 返回更新的数据slice, 如果为空则表示没有新的成绩
// 操作步骤: 抓取用户成绩(N) -> 获取数据用户已有成绩(M) -> 数据库删除数据(M-N) -> 数据库插入数据(N-M)
// 用户成绩更新大多数情况都是新增，对于只是更新个别字段的成绩直接删除即可
func UpdateGrades(userID uint) ([]Grade, error) {
	// 获取教务处句柄
	c, err := GetJwc(userID)
	if err != nil {
		return nil, err
	}

	// 获取全部成绩
	grades := grade.GetALL(c)
	if len(grades) == 0 {
		return nil, errors.New("没有从教务处获取到成绩信息")
	}
	// 获取教务处句柄
	c, err = GetJwc(userID)
	if err != nil {
		return nil, err
	}
	// 获取不及格成绩
	failGrades := grade.GetNotPass(c)
	grades = append(grades, failGrades...)

	// slice to set
	newGradeSet := mapset.NewSet()
	for _, g := range grades {
		newGradeSet.Add(Grade{
			CourseID:   g.CourseID,
			LessonID:   g.LessonID,
			CourseName: g.CourseName,
			Credit:     g.Credit,
			CourseType: g.CourseType,
			Grade:      g.Grade,
			TermName:   g.TermName,
			Year:       g.Year,
			Term:       g.Term,
		})
	}

	// 从数据库取出现有数据
	var oldGrades []Grade
	err = DB().Where("user_id = ?", userID).Find(&oldGrades).Error
	if err != nil {
		return nil, err
	}
	// slice to set
	oldGradeSet := mapset.NewSet()
	oldGradeIDs := map[string]uint{}
	for _, v := range oldGrades {
		oldGradeIDs[fmt.Sprintf("%s-%s", v.CourseID, v.LessonID)] = v.ID
		v.UserID = 0
		v.Model = Model{}
		oldGradeSet.Add(v)
	}

	// 新建事务，准备开始数据库操作
	tx := DB().Begin()

	// 删除需要删除的数据
	deleteSet := oldGradeSet.Difference(newGradeSet)
	it := deleteSet.Iterator()
	var deleteIDs []uint
	for elem := range it.C {
		g := elem.(Grade)
		if id, ok := oldGradeIDs[fmt.Sprintf("%s-%s", g.CourseID, g.LessonID)]; ok {
			deleteIDs = append(deleteIDs, id)
		}
	}
	if len(deleteIDs) != 0 {
		tx.Unscoped().Where(deleteIDs).Delete(Grade{})
	}

	// 新增更新的数据
	createSet := newGradeSet.Difference(oldGradeSet)
	it = createSet.Iterator()
	var updateGrades []Grade
	for elem := range it.C {
		g := elem.(Grade)
		g.UserID = userID
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
