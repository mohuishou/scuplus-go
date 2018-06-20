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
