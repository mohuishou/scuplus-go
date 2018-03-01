package tasks

import (
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
