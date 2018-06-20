package tasks

import (
	"time"

	"github.com/mohuishou/scuplus-go/model"
)

// CalCourse 计算课程
func CalCourse(cid uint) error {
	var course model.CourseCount
	if err := model.DB().Find(&course, cid).Error; err == nil {
		// 获取所有的课程号数据
		countGradeAll(&course)
		calCourseGrade(&course)
	}
	return nil
}

// countGradeAll 统计所有成绩数据
func countGradeAll(courseCount *model.CourseCount) {
	// 计算该门课程的数据
	var (
		all float64
		sum struct {
			Total float64
		}
		fail float64
	)
	// 成绩信息统计
	model.DB().Model(&model.Grade{}).Where(model.Grade{
		CourseID: courseCount.CourseID,
		LessonID: courseCount.LessonID,
	}).Count(&all)
	model.DB().Model(&model.Grade{}).Where(model.Grade{
		CourseID: courseCount.CourseID,
		LessonID: courseCount.LessonID,
	}).Where("grade > ? and grade < ?", 0, 60).Count(&fail)
	model.DB().Model(&model.Grade{}).Where(model.Grade{
		CourseID: courseCount.CourseID,
		LessonID: courseCount.LessonID,
	}).Select("sum(grade) total").Scan(&sum)
	// 计算平均分
	model.DB().Where("course_id = ? and lesson_id = ?", courseCount.CourseID, courseCount.LessonID).FirstOrCreate(courseCount)
	// 数据更新
	if all > 0 {
		courseCount.AvgGrade = sum.Total / all
		courseCount.FailRate = fail / all
	} else {
		courseCount.AvgGrade = 0
		courseCount.FailRate = 0
	}
	courseCount.GradeAll = int(all)
	model.DB().Save(&courseCount)
}

// calCourseGrade 统计历史成绩数据
func calCourseGrade(courseCount *model.CourseCount) {
	// 开始年份
	startYear := 2014
	// 获取当前年份
	endYear := time.Now().Year()
	// 从数据库统计相关人数
	for i := startYear; i <= endYear; i++ {
		cg := model.CourseGrade{
			CourseID: courseCount.CourseID,
			LessonID: courseCount.LessonID,
			Year:     i,
		}
		model.DB().Where(cg).FirstOrCreate(&cg)
		// 0 ~ 60
		model.DB().Model(&model.Grade{}).Where(model.Grade{
			Year:     i,
			CourseID: courseCount.CourseID,
			LessonID: courseCount.LessonID,
		}).Where("grade > ? and grade < ?", 0, 60).Count(&cg.G0)
		// 60 ~ 70
		model.DB().Model(&model.Grade{}).Where(model.Grade{
			Year:     i,
			CourseID: courseCount.CourseID,
			LessonID: courseCount.LessonID,
		}).Where("grade > ? and grade < ?", 60, 70).Count(&cg.G60)
		// 70 ~ 80
		model.DB().Model(&model.Grade{}).Where(model.Grade{
			Year:     i,
			CourseID: courseCount.CourseID,
			LessonID: courseCount.LessonID,
		}).Where("grade > ? and grade < ?", 70, 80).Count(&cg.G70)
		// 80 ~ 90
		model.DB().Model(&model.Grade{}).Where(model.Grade{
			Year:     i,
			CourseID: courseCount.CourseID,
			LessonID: courseCount.LessonID,
		}).Where("grade > ? and grade < ?", 80, 90).Count(&cg.G80)
		// 90 ~ 100
		model.DB().Model(&model.Grade{}).Where(model.Grade{
			Year:     i,
			CourseID: courseCount.CourseID,
			LessonID: courseCount.LessonID,
		}).Where("grade > ? and grade < ?", 90, 100).Count(&cg.G90)

		// 计算是否每一项都为0
		if cg.G0 != 0 || cg.G60 != 0 || cg.G70 != 0 || cg.G80 != 0 || cg.G90 != 0 {
			model.DB().Save(&cg)
		}
	}
}
