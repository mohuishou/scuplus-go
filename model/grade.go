package model

import "log"

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
	TermName   string `json:"term_name"`
}

// GetGrades 获取用户的所有成绩
func GetGrades(userID uint) []Grade {
	grades := make([]Grade, 0)
	if err := DB().Where("user_id", userID).Find(&grades); err != nil {
		log.Printf("[Error] GetGrades Fail, userID: %d", userID)
	}
	return grades
}

// UpdateGrades 更新用户的所有成绩
func UpdateGrades(userID uint) error {
	// 获取教务处句柄
	jwc, err := GetJwc(userID)
	if err != nil {
		return err
	}

	// 获取全部成绩
	grades, err := jwc.GPAAll()
	if err != nil {
		return err
	}

	tx := DB().Begin()

	for _, gradeArr := range grades {
		for _, g := range gradeArr {
			grade := Grade{
				CourseID:   g.CourseID,
				LessonID:   g.LessonID,
				CourseName: g.CourseName,
				Credit:     g.Credit,
				CourseType: g.CourseType,
				Grade:      g.Grade,
				Term:       g.Term,
				TermName:   g.TermName,
			}

			oldGrade := Grade{}
			tx.Where(Grade{
				UserID:   userID,
				CourseID: g.CourseID,
				LessonID: g.LessonID,
			}).Find(&oldGrade)

			// 有则更新，无则创建
			if oldGrade.ID != 0 {
				if err := tx.Where("id", oldGrade.ID).Update(grade).Error; err != nil {
					tx.Rollback()
					log.Println("[Error]: UpdateGrades", grade)
					return err
				}
			} else {
				if err := tx.Create(&grade).Error; err != nil {
					tx.Rollback()
					log.Println("[Error]: UpdateGrades", grade)
					return err
				}
			}

			tx.Commit()
		}
	}
	return nil
}
