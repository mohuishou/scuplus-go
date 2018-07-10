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
	scope := model.DB().Where(model.Grade{
		CourseID: courseCount.CourseID,
		LessonID: courseCount.LessonID,
	})
	gradeScope := scope.Model(&model.Grade{})
	gradeScope.Count(&all)
	gradeScope.Count(&all).Where("grade > ? and grade < ?", 0, 60).Count(&fail)
	gradeScope.Count(&all).Select("sum(grade) total").Scan(&sum)

	// 计算平均分
	scope.FirstOrCreate(courseCount)
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
	// 统计不及格人数
	// year = -1 尚不及格
	countGrades(courseCount, -1)
	// year = -2 曾不及格
	countGrades(courseCount, -2)

	// 开始年份
	startYear := 2014
	// 获取当前年份
	endYear := time.Now().Year()
	// 从数据库统计相关人数
	for i := startYear; i <= endYear; i++ {
		countGrades(courseCount, i)
	}
}

func countGrades(courseCount *model.CourseCount, year int) {
	cg := model.CourseGrade{
		CourseID: courseCount.CourseID,
		LessonID: courseCount.LessonID,
		Year:     year,
	}
	model.DB().Where(cg).FirstOrCreate(&cg)

	// 60 ~ 100 分统计， 不及格成绩的年份为-1 或者 -2
	if year > 0 {
		for i := 60; i < 100; i += 10 {
			countGradesRange(courseCount, &cg, year, i)
		}
	} else {
		countGradesRange(courseCount, &cg, year, 0)
	}

	// 计算是否每一项都为0
	if (cg.G0 + cg.G60 + cg.G70 + cg.G80 + cg.G90) != 0 {
		model.DB().Save(&cg)
	}
}

func countGradesRange(cc *model.CourseCount, cg *model.CourseGrade, year, gradeRange int) {
	min, max := gradeRange, gradeRange+10
	if min == 0 {
		max = 60
	}
	scope := model.DB().Model(&model.Grade{}).Where(model.Grade{
		Year:     year,
		CourseID: cc.CourseID,
		LessonID: cc.LessonID,
	}).Where("grade > ? and grade < ?", min, max)

	switch gradeRange {
	case 0:
		scope.Count(&cg.G0)
	case 60:
		scope.Count(&cg.G60)
	case 70:
		scope.Count(&cg.G70)
	case 80:
		scope.Count(&cg.G80)
	case 90:
		scope.Count(&cg.G90)
	}
}
