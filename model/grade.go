package model

import (
	"errors"
	"log"

	"github.com/mohuishou/scu/jwc/grade"
)

// Grade 成绩
type Grade struct {
	Model
	UserID     uint   `json:"user_id"`
	CourseID   string `json:"course_id"`
	LessonID   string `json:"lesson_id"`
	CourseName string `json:"course_name"`
	Credit     string `json:"credit"`
	CourseType string `json:"course_type"`
	Grade      string `json:"grade"`
	Term       int    `json:"term"`
	Year       int    `json:"year"`
	TermName   string `json:"term_name"`
}

// GetGrades 获取用户的所有成绩
func GetGrades(userID uint) []Grade {
	grades := make([]Grade, 0)
	if err := DB().Where("user_id = ?", userID).Find(&grades).Error; err != nil {
		log.Printf("[Error] GetGrades Fail, userID: %d, err: %s", userID, err.Error())
	}
	return grades
}

// UpdateGrades 更新用户的所有成绩
func UpdateGrades(userID uint) error {
	// 获取教务处句柄
	c, err := GetJwc(userID)
	if err != nil {
		return err
	}

	// 获取全部成绩
	grades := grade.GetALL(c)
	if len(grades) == 0 {
		return errors.New("没有从教务处获取到成绩信息")
	}

	tx := DB().Begin()

	for _, g := range grades {
		grade := Grade{
			UserID:     userID,
			CourseID:   g.CourseID,
			LessonID:   g.LessonID,
			CourseName: g.CourseName,
			Credit:     g.Credit,
			CourseType: g.CourseType,
			Grade:      g.Grade,
			Term:       g.Term,
			Year:       g.Year,
			TermName:   g.TermName,
		}

		oldGrade := Grade{}
		tx.Where(Grade{
			UserID:   userID,
			CourseID: g.CourseID,
			LessonID: g.LessonID,
			Term:     g.Term,
		}).Find(&oldGrade)

		// 有则更新，无则创建
		if oldGrade.ID != 0 {
			if err := tx.Model(&oldGrade).Updates(grade).Error; err != nil {
				tx.Rollback()
				log.Println("[Error]: UpdateGrades", err, grade)
				return err
			}
		} else {
			if err := tx.Create(&grade).Error; err != nil {
				tx.Rollback()
				log.Println("[Error]: UpdateGrades", err, grade)
				return err
			}
		}
	}
	tx.Commit()

	return nil
}
