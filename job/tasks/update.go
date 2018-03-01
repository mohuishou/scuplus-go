package tasks

import (
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/mohuishou/scuplus-go/job"
	"github.com/mohuishou/scuplus-go/model"
)

// UpdateAll 更新用户相关信息，包括但不限于
// 教务处相关: 成绩/考表
// 图书馆相关: 借阅信息
// 一卡通相关: 交易流水
func UpdateAll(uid uint) error {
	// 更新成绩
	updateGrades, err := model.UpdateGrades(uid)
	if err != nil {
		return err
	}
	if len(updateGrades) > 0 {
		// 有新的成绩通知，添加到通知队列
		sign := &tasks.Signature{
			Name: "notify_grade",
			Args: []tasks.Arg{
				{
					Type:  "uint",
					Value: uid,
				},
				{
					Type:  "string",
					Value: updateGrades[0].CourseName,
				},
				{
					Type:  "string",
					Value: updateGrades[0].Grade,
				},
				{
					Type:  "int",
					Value: updateGrades[0].Credit,
				},
				{
					Type:  "int",
					Value: len(updateGrades),
				},
			},
		}
		server, err := job.StartServer()
		if err != nil {
			return err
		}
		_, err = server.SendTask(sign)
		if err != nil {
			return err
		}
	}

	// 更新考表
	err = model.UpdateExam(uid)
	if err != nil {
		return err
	}

	// 更新借阅信息
	_, err = model.UpdateLibraryBook(uid, 0)
	if err != nil {
		return err
	}

	// 交易流水
	err = model.UpdateEcard(uid)
	if err != nil {
		return err
	}

	return nil
}
