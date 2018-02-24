package tasks

import (
	"github.com/mohuishou/scuplus-go/model"
)

// UpdateAll 更新用户相关信息，包括但不限于
// 教务处相关: 课表/成绩/考表
// 图书馆相关: 借阅信息
// 一卡通相关: 交易流水
func UpdateAll(uid uint) (err error) {
	// 更新课表

	// 更新成绩
	err = model.UpdateGrades(uid)
	if err != nil {
		return
	}

	// 更新考表
	err = model.UpdateExam(uid)
	if err != nil {
		return
	}

	// 更新借阅信息
	_, err = model.UpdateLibraryBook(uid, 0)
	if err != nil {
		return
	}

	// 交易流水
	err = model.UpdateEcard(uid)
	if err != nil {
		return
	}

	return nil
}
