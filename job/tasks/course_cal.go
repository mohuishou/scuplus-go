package tasks

import (
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	"log"

	"github.com/mohuishou/scuplus-go/model"
)

// CalCourse 计算课程
func CalCourse(cid uint) error {
	var course model.Course
	if err := model.DB().Find(&course, cid).Error; err == nil {
		// 获取所有的课程号数据
		calCourceAll(course)
		calCourceGrade(course)
	}
	return nil
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
	// 成绩信息统计
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

	courseCount := model.CourseCount{
		CourseID: course.CourseID,
		LessonID: course.LessonID,
		Name:     course.Name,
		Day:      course.Day,
		Credit:   course.Credit,
		Campus:   course.Campus,
	}
	// 计算平均分
	model.DB().Where("course_id = ? and lesson_id = ?", course.CourseID, course.LessonID).FirstOrCreate(&courseCount)
	if all > 0 {
		courseCount.AvgGrade = sum.Total / all
		courseCount.FailRate = fail / all
	} else {
		courseCount.AvgGrade = 0
		courseCount.FailRate = 0
	}

	courseCount.GradeAll = int(all)
	// 评教信息统计
	model.DB().Model(&model.CourseEvaluate{}).Where("course_id = ? and lesson_id = ? and status = 1", course.CourseID, course.LessonID).Select([]string{"AVG(star) star"}).Scan(&courseCount)
	countEva(course, "call_name").Scan(&courseCount)
	countEva(course, "exam_type").Scan(&courseCount)
	countEva(course, "task").Scan(&courseCount)

	countEvaStar(course, "good", "3").Scan(&courseCount)
	countEvaStar(course, "normal", "2").Scan(&courseCount)
	countEvaStar(course, "bad", "1").Scan(&courseCount)

	// 教师统计
	teachers := []model.Teacher{}
	model.DB().Model(&course).Related(&teachers, "Teachers")
	courseCount.Teacher = ""
	for _, teacher := range teachers {
		courseCount.Teacher = courseCount.Teacher + "," + teacher.Name
	}
	courseCount.Teacher = strings.Trim(courseCount.Teacher, ",")
	log.Println("即将更新课程数据：", courseCount)
	model.DB().Save(&courseCount)
}

func countEva(course model.Course, name string) *gorm.DB {
	return model.DB().Model(&model.CourseEvaluate{}).Where("course_id = ? and lesson_id = ? and status = 1", course.CourseID, course.LessonID).Select([]string{name, "Count(*) c"}).Group(name).Order("c desc").Limit(1)
}

func countEvaStar(course model.Course, name, star string) *gorm.DB {
	return model.DB().Model(&model.CourseEvaluate{}).Where("course_id = ? and lesson_id = ? and star = ? and status = 1", course.CourseID, course.LessonID, star).Select([]string{"Count(*) " + name}).Group("star")
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
