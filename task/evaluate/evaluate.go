package main

import (
	"github.com/mohuishou/scuplus-go/model"
)

// 课程评教更新
func main() {

	// 获取已评教的信息，根据课程号与课序号添加课程名
	courseEvaList := []model.CourseEvaluate{}
	model.DB().Find(&courseEvaList)
	for _, v := range courseEvaList {
		if v.CourseName != "" {
			continue
		}
		c := model.Course{}
		model.DB().Where("course_id = ? and lesson_id = ? ",
			v.CourseID,
			v.LessonID,
		).First(&c)
		model.DB().Model(&v).Update("course_name", c.Name)
	}

	grades := []model.Grade{}
	model.DB().Select([]string{
		"user_id", "course_id", "lesson_id", "course_name",
	}).Find(&grades)
	for _, v := range grades {
		model.CBNewCourseEvaluate(v.UserID, v.CourseID, v.LessonID, v.CourseName)
	}

	schs := []model.Schedule{}
	model.DB().Select([]string{
		"user_id", "course_id", "lesson_id", "course_name",
	}).Find(&schs)
	for _, v := range schs {
		model.CBNewCourseEvaluate(v.UserID, v.CourseID, v.LessonID, v.CourseName)
	}
}
