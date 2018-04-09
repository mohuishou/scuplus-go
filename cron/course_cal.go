package main

import (
	"time"

	"github.com/mohuishou/scuplus-go/model"
)

// 计算课程分数
func calCourse() {
	// 获取所有的课程号数据
	courses := []model.Course{}
	model.DB().Find(&courses)
	for _, course := range courses {
		calCourceAll(course)
		calCourceGrade(course)
	}
}

func calCourceAll(course model.Course) {
	// 计算该门课程的数据
	var (
		all float64
		sum struct {
			Total float64
		}
		fail float64
	)
	model.DB().Model(&model.Grade{}).Where(model.Grade{
		CourseID: course.CourseID,
		LessonID: course.LessonID,
	}).Count(&all)
	model.DB().Model(&model.Grade{}).Where(model.Grade{
		CourseID: course.CourseID,
		LessonID: course.LessonID,
	}).Where("grade > ? and grade < ?", 0, 60).Count(&fail)

	model.DB().Model(&model.Grade{}).Where(model.Grade{
		CourseID: course.CourseID,
		LessonID: course.LessonID,
	}).Select("sum(grade) total").Scan(&sum)
	// 计算平均分
	course.AvgGrade = sum.Total / all
	course.FailRate = fail / all
	model.DB().Save(&course)
}

func calCourceGrade(course model.Course) {
	// 开始年份
	startYear := 2014
	// 获取当前年份
	endYear := time.Now().Year()
	// 从数据库统计相关人数
	for i := startYear; i <= endYear; i++ {
		cg := model.CourseGrade{
			CourseID: course.CourseID,
			LessonID: course.LessonID,
			Year:     i,
		}
		model.DB().Where(cg).FirstOrCreate(&cg)
		// 0 ~ 60
		model.DB().Model(&model.Grade{}).Where(model.Grade{
			Year:     i,
			CourseID: course.CourseID,
			LessonID: course.LessonID,
		}).Where("grade > ? and grade < ?", 0, 60).Count(&cg.G0)
		// 60 ~ 70
		model.DB().Model(&model.Grade{}).Where(model.Grade{
			Year:     i,
			CourseID: course.CourseID,
			LessonID: course.LessonID,
		}).Where("grade > ? and grade < ?", 60, 70).Count(&cg.G60)
		// 70 ~ 80
		model.DB().Model(&model.Grade{}).Where(model.Grade{
			Year:     i,
			CourseID: course.CourseID,
			LessonID: course.LessonID,
		}).Where("grade > ? and grade < ?", 70, 80).Count(&cg.G70)
		// 80 ~ 90
		model.DB().Model(&model.Grade{}).Where(model.Grade{
			Year:     i,
			CourseID: course.CourseID,
			LessonID: course.LessonID,
		}).Where("grade > ? and grade < ?", 80, 90).Count(&cg.G80)
		// 90 ~ 100
		model.DB().Model(&model.Grade{}).Where(model.Grade{
			Year:     i,
			CourseID: course.CourseID,
			LessonID: course.LessonID,
		}).Where("grade > ? and grade < ?", 90, 100).Count(&cg.G90)

		// 计算是否每一项都为0
		if cg.G0 != 0 || cg.G60 != 0 || cg.G70 != 0 || cg.G80 != 0 || cg.G90 != 0 {
			model.DB().Save(&cg)
		}
	}
}
