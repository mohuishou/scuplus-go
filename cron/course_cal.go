package main

import (
	"log"

	"github.com/RichardKnop/machinery/v1/tasks"

	"github.com/mohuishou/scuplus-go/job"
	"github.com/mohuishou/scuplus-go/model"
)

// 计算课程分数
func calCourse() {
	courses := []model.CourseCount{}
	model.DB().Find(&courses)
	for _, course := range courses {
		sign := &tasks.Signature{
			Name: "course_count",
			Args: []tasks.Arg{
				{
					Type:  "uint",
					Value: course.ID,
				},
			},
		}

		_, err := job.Server.SendTask(sign)
		if err != nil {
			log.Println("cron error course_count", err)
		}
	}
}

// courseCountEvaluate 课程评价统计
func courseCountEvaluate() {
	var courseEvas []model.CourseEvaluate
	model.DB().Where("status = 1").Order("course_id desc, lesson_id desc").Find(&courseEvas)
	var (
		courseID      string
		lessonID      string
		all           int
		good          int
		normal        int
		bad           int
		star          float64
		sum           int
		callNameCount map[int]int
		taskCount     map[int]int
		examTypeCount map[int]int
	)
	for _, courseEva := range courseEvas {
		if courseID != courseEva.CourseID || lessonID != courseEva.LessonID {
			if courseID != "" && lessonID != "" {
				star = float64(sum) / float64(all)
				courseCount := model.CourseCount{}
				if err := model.DB().Where("course_id = ? and lesson_id = ?", courseID, lessonID).First(&courseCount).Error; err != nil {
					return
				}
				courseCount.Star = star
				courseCount.Good = good
				courseCount.Normal = normal
				courseCount.Bad = bad
				courseCount.ExamType = max(examTypeCount)
				courseCount.Task = max(taskCount)
				courseCount.CallName = max(callNameCount)
				model.DB().Save(&courseCount)
			}
			courseID = courseEva.CourseID
			lessonID = courseEva.LessonID
			// 初始化
			all, sum, good, normal, bad = 0, 0, 0, 0, 0
			// 统计
			taskCount = make(map[int]int)
			callNameCount = make(map[int]int)
			examTypeCount = make(map[int]int)
		}
		all++
		sum = sum + courseEva.Star
		taskCount[courseEva.Task]++
		callNameCount[courseEva.CallName]++
		examTypeCount[courseEva.ExamType]++
		switch courseEva.Star {
		case 3:
			good++
		case 2:
			normal++
		case 1:
			bad++
		}
	}
}

func max(data map[int]int) int {
	maxVal := 0
	maxKey := 0
	for k, v := range data {
		if maxVal < v {
			maxKey = k
			maxVal = v
		}
	}
	return maxKey
}
