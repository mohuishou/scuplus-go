package model

import "testing"

func TestCourse_AfterSave(t *testing.T) {
	DB().Create(&Course{
		LessonID: "48",
		CourseID: "-234234242",
		Name:     "测试课程",
	})
}
