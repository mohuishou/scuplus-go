package main

import (
	"log"
	"time"

	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/mohuishou/scuplus-go/job"
	"github.com/mohuishou/scuplus-go/model"
)

const pageSize = 1000

func main() {

}

func updateAll() {
	count := 0
	model.DB().Table("users").Where("verify = 1").Count(&count)
	for i := 0; i < (count/pageSize + 1); i++ {
		users := []model.User{}
		model.DB().Select([]string{"id"}).Offset((i - 1) * pageSize).Limit(pageSize).Find(&users)
		for _, user := range users {
			sign := &tasks.Signature{
				Name: "update_all",
				Args: []tasks.Arg{
					{
						Type:  "uint",
						Value: user.ID,
					},
				},
			}

			_, err := job.Server.SendTask(sign)
			if err != nil {
				log.Println("cron error update all", err)
			}
		}
	}
}

func book() {
	count := 0
	model.DB().Table("users").Where("verify = 1").Count(&count)
	for i := 0; i < (count/pageSize + 1); i++ {
		users := []model.User{}
		model.DB().Select([]string{"id"}).Offset((i - 1) * pageSize).Limit(pageSize).Find(&users)
		for _, user := range users {

			// 获取即将到期的书籍

			sign := &tasks.Signature{
				Name: "update_all",
				Args: []tasks.Arg{
					{
						Type:  "uint",
						Value: user.ID,
					},
				},
			}

			_, err := job.Server.SendTask(sign)
			if err != nil {
				log.Println("cron error update all", err)
			}
		}
	}
}

func exam() {
	count := 0
	model.DB().Table("users").Where("verify = 1").Count(&count)
	for i := 0; i < (count/pageSize + 1); i++ {
		users := []model.User{}
		model.DB().Select([]string{"id"}).Offset((i - 1) * pageSize).Limit(pageSize).Find(&users)
		for _, user := range users {
			// 获取最近的考试
			now := time.Now()
			var exam model.Exam
			model.DB().Where("user_id = ? and start_time > ?", user.ID, now).Find(&exam)
			day := (exam.StartTime.Unix() - now.Unix()) / 3600 / 24
			if day == 30 || day == 7 || day == 0 {
				sign := &tasks.Signature{
					Name: "notify_exam",
					Args: []tasks.Arg{
						{
							Type:  "uint",
							Value: user.ID,
						},
						{
							Type:  "string",
							Value: exam.CourseName,
						},
						{
							Type:  "string",
							Value: exam.Date,
						},
						{
							Type:  "string",
							Value: exam.Time,
						},
						{
							Type:  "string",
							Value: exam.Campus + exam.Building + exam.Classroom,
						},
						{
							Type:  "string",
							Value: exam.Site,
						},
						{
							Type:  "string",
							Value: exam.Name,
						},
						{
							Type:  "int64",
							Value: day,
						},
					},
				}
				_, err := job.Server.SendTask(sign)
				if err != nil {
					log.Println("cron error exam", err)
				}
			}

		}
	}
}
