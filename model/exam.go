package model

import (
	"github.com/mohuishou/scu/jwc/exam"
)

// Exam 考表
type Exam struct {
	Model
	Name       string `json:"name"`
	Campus     string `json:"campus"`
	Building   string `json:"building"`
	Classroom  string `json:"classroom"`
	CourseName string `json:"course_name"`
	Week       string `json:"week"`
	Day        string `json:"day"`
	Date       string `json:"date"`
	Time       string `json:"time"`
	Site       string `json:"site"`
	Card       string `json:"card"`
	Comment    string `json:"comment"`
}

// convertExam 转换为model
func convertExam(e exam.Exam) Exam {
	return Exam{}
}
