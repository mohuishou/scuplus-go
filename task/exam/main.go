package main

import (
	"time"

	"log"

	"sync"

	"strings"

	"github.com/mohuishou/scuplus-go/model"
)

// fixed exam start_time
func main() {
	e := model.Exam{}
	model.DB().Find(&e, 16)
	fixedOne(&e)
	w := sync.WaitGroup{}
	count := 0
	page, pageSize := 1, 100
	model.DB().Model(&model.Exam{}).Count(&count)
	log.Println("Exam总计", count, "条")
	for ; page <= ((count / pageSize) + 1); page++ {
		w.Add(1)
		go func(page, pageSize int) {
			log.Println("第", page, "页开始，总计", page, "页")
			scope := model.DB().Offset((page - 1) * pageSize).Limit(pageSize)
			exams := []model.Exam{}
			scope.Select([]string{
				"id", "start_time", "date", "time",
			}).Find(&exams)
			for _, v := range exams {
				fixedOne(&v)
			}
			log.Println("第", page, "页结束，总计", page, "页")
			w.Done()
		}(page, pageSize)
	}
	w.Wait()
}

func fixedOne(exam *model.Exam) {
	timeStr := exam.Date + " " + strings.Split(exam.Time, "-")[0]
	startTime, err := time.ParseInLocation("2006-01-02 15:04", timeStr, time.Local)
	if err != nil {
		log.Println("error", err)
		return
	}
	err = model.DB().Model(exam).Update(
		"start_time",
		startTime).Error
	if err != nil {
		log.Println("error", err)
	}

}
