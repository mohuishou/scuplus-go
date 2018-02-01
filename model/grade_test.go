package model

import (
	"testing"
)

func TestUpdateGrades(t *testing.T) {
	g1 := Grade{
		UserID:     uint(1),
		CourseID:   "123456",
		LessonID:   "01",
		CourseName: "xx",
		Credit:     "xx",
		CourseType: "xx",
		Grade:      "xx",
		Term:       2,
		Year:       2017,
		TermName:   "xx",
	}

	g2 := Grade{
		UserID:     uint(1),
		CourseID:   "123456",
		LessonID:   "01",
		CourseName: "xx",
		Credit:     "xx",
		CourseType: "xx",
		Grade:      "xx",
		Term:       2,
		Year:       2017,
		TermName:   "xx",
	}

	if g1 == g2 {
		t.Log("kkkkkkkkk")
	}
}
