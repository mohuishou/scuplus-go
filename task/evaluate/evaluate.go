package main

import (
	"log"

	"sync"

	"github.com/mohuishou/scuplus-go/model"
)

// 课程评教更新
func main() {

	// 获取已评教的信息，根据课程号与课序号添加课程名
	//courseEvaList := []model.CourseEvaluate{}
	//model.DB().Find(&courseEvaList)
	//for _, v := range courseEvaList {
	//	if v.CourseName != "" {
	//		continue
	//	}
	//	c := model.Course{}
	//	model.DB().Where("course_id = ? and lesson_id = ? ",
	//		v.CourseID,
	//		v.LessonID,
	//	).First(&c)
	//	model.DB().Model(&v).Update("course_name", c.Name)
	//}

	updateGrade()
	updateSch()
}

func updateGrade() {
	w := sync.WaitGroup{}
	count := 0
	page, pageSize := 1, 200
	model.DB().Model(&model.Grade{}).Count(&count)
	log.Println("成绩总计", count, "条")
	for ; page < ((count / pageSize) + 1); page++ {
		w.Add(1)
		go func(page, pageSize int) {
			log.Println("第", page, "页开始，总计", page, "页")
			scope := model.DB().Offset((page - 1) * pageSize).Limit(pageSize)
			grades := []model.Grade{}
			scope.Select([]string{
				"user_id", "course_id", "lesson_id", "course_name",
			}).Find(&grades)
			for _, v := range grades {
				model.CBNewCourseEvaluate(v.UserID, v.CourseID, v.LessonID, v.CourseName)
			}
			log.Println("第", page, "页结束，总计", page, "页")
			w.Done()
		}(page, pageSize)
	}
	w.Wait()
}

func updateSch() {
	w := sync.WaitGroup{}
	count := 0
	page, pageSize := 1, 200
	model.DB().Model(&model.Schedule{}).Count(&count)
	log.Println("课程表总计", count, "条")
	for ; page < ((count / pageSize) + 1); page++ {
		w.Add(1)
		go func(page, pageSize int) {
			log.Println("第", page, "页开始，总计", page, "页")
			scope := model.DB().Offset((page - 1) * pageSize).Limit(pageSize)
			schs := []model.Schedule{}
			scope.Select([]string{
				"user_id", "course_id", "lesson_id", "course_name",
			}).Find(&schs)
			for _, v := range schs {
				model.CBNewCourseEvaluate(v.UserID, v.CourseID, v.LessonID, v.CourseName)
			}
			log.Println("第", page, "页结束，总计", page, "页")
			w.Done()
		}(page, pageSize)
	}
	w.Wait()
}
