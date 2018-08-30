package model

import (
	"errors"
	"log"

	"fmt"

	"github.com/mohuishou/scu/jwc"
	"github.com/mohuishou/scu/jwc/grade"
)

// Grade 成绩
// Note: 唯一性索引添加上grade，因为在不及格成绩当中可能会出现同一个学期有两次不及格成绩的情况
// 风险: 有导致其他成绩出现两个不同成绩的风险
type Grade struct {
	Model
	UserID     uint    `json:"user_id" gorm:"unique_index:u_grade"`
	CourseID   string  `json:"course_id" gorm:"unique_index:u_grade"`
	LessonID   string  `json:"lesson_id" gorm:"unique_index:u_grade"`
	CourseName string  `json:"course_name"`
	Credit     string  `json:"credit"`
	CourseType string  `json:"course_type"`
	Grade      float64 `json:"grade" gorm:"unique_index:u_grade"`
	GPA        float64 `json:"gpa"`
	Term       int     `json:"term" gorm:"unique_index:u_grade"` //0: 秋季学期, 1: 春季学期
	Year       int     `json:"year" gorm:"unique_index:u_grade"`
	TermName   string  `json:"term_name"`
}

type Grades []Grade

// getKey 一门成绩的唯一标识
func (g Grade) getKey() string {
	return fmt.Sprintf("%d-%s-%s-%f-%d-%d",
		g.UserID,
		g.CourseID,
		g.LessonID,
		g.Grade,
		g.Year,
		g.Term,
	)
}

// getKeyMap 获取一组成绩的唯一标识集合
func (grades Grades) getKeyMap() map[string]bool {
	keyMap := make(map[string]bool)
	for _, g := range grades {
		keyMap[g.getKey()] = true
	}
	return keyMap
}

// Difference result = grades - others
func (grades Grades) Difference(other Grades) (Grades, []uint) {
	keyMap := other.getKeyMap()
	differenceGrades := make(Grades, 0)
	ids := make([]uint, 0)
	for _, g := range grades {
		if _, ok := keyMap[g.getKey()]; !ok {
			differenceGrades = append(differenceGrades, g)
			if g.ID != 0 {
				ids = append(ids, g.ID)
			}
		}
	}
	return differenceGrades, ids
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
	// 从教务处获取成绩
	grades, err := getAllGradesFromJwc(userID)
	if err != nil {
		return nil, err
	}
	log.Printf("%d 从教务处获取成绩: %d条", userID, len(grades))

	// 从数据库取出现有数据
	var oldGrades Grades
	DB().Where("user_id = ?", userID).Find(&oldGrades)
	log.Printf("%d 已有成绩: %d条", userID, len(oldGrades))

	// 获取需要删除的ids
	_, ids := oldGrades.Difference(grades)
	if len(ids) != 0 {
		if err := DB().Unscoped().Where(ids).Delete(Grade{}).Error; err != nil {
			log.Println("grade delete err:", err)
			return nil, err
		}
	}
	log.Printf("%d 删除成绩: %d条", userID, len(ids))

	// 新增更新的数据
	updateGrades, _ := grades.Difference(oldGrades)
	log.Printf("%d 需要新增成绩: %d条", userID, len(updateGrades))

	updates := make(Grades, 0)
	for _, g := range updateGrades {
		if err := DB().Create(&g).Error; err != nil {
			log.Println("grade create err:", g, err)
		} else {
			updates = append(updates, g)
		}
	}
	return updates, nil
}

// getAllGradesFromJwc 从教务处获取所有成绩
func getAllGradesFromJwc(userID uint) (Grades, error) {
	// 获取教务处句柄
	c, err := GetJwc(userID)
	if err != nil {
		return nil, err
	}

	// 获取全部成绩
	grades, err := grade.GetALL(c)
	if err != nil {
		return nil, err
	}

	// 获取教务处句柄
	c, err = GetJwc(userID)
	if err != nil {
		return nil, err
	}

	defer jwc.Logout(c)

	// 获取不及格成绩
	failGrades, err := grade.GetNotPass(c)
	if err != nil {
		return nil, err
	}
	grades = append(grades, failGrades...)
	if len(grades) == 0 {
		return nil, errors.New("没有从教务处获取到成绩信息")
	}

	// 类型转换
	newGrades := make([]Grade, len(grades))
	for i, g := range grades {
		newGrades[i] = convertGrade(userID, g)
	}
	return newGrades, nil
}

func convertGrade(uid uint, g grade.Grade) Grade {
	return Grade{
		UserID:     uid,
		CourseID:   g.CourseID,
		LessonID:   g.LessonID,
		CourseName: g.CourseName,
		Credit:     g.Credit,
		CourseType: g.CourseType,
		Grade:      g.Grade,
		TermName:   g.TermName,
		Year:       g.Year,
		GPA:        g.GPA,
		Term:       g.Term,
	}
}
