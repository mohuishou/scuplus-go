package tasks

import (
	"log"

	"strings"

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
		return errorHandle("jwc", uid, err)
	}
	if len(updateGrades) > 0 {
		// 有新的成绩通知，添加到通知队列
		err = NotifyGrade(uid, updateGrades[0].CourseName, updateGrades[0].Grade, updateGrades[0].Credit, len(updateGrades))
		log.Println("notify error", err)
	}

	// 更新考表
	err = model.UpdateExam(uid)
	if err != nil {
		return errorHandle("jwc", uid, err)
	}

	// 更新借阅信息
	_, err = model.UpdateLibraryBook(uid, 0)
	if err != nil {
		return errorHandle("library", uid, err)
	}

	// 交易流水
	err = model.UpdateEcard(uid)
	if err != nil {
		return errorHandle("my", uid, err)
	}

	return nil
}

// UpdateForNew 新用户
func UpdateForNew(uid uint) error {
	// 更新成绩
	_, err := model.UpdateGrades(uid)
	if err != nil {
		log.Println(err)
	}

	// 更新考表
	err = model.UpdateExam(uid)
	if err != nil {
		log.Println(err)
	}

	// 交易流水
	err = model.UpdateEcard(uid)
	if err != nil {
		log.Println(err)
	}
	return err
}

func errorHandle(verifyType string, uid uint, err error) error {
	log.Println(verifyType, err)
	if !strings.Contains(err.Error(), "密码") {
		return err
	}
	user := model.DB().Model(&model.User{
		Model: model.Model{ID: uid},
	})
	userLib := model.DB().Model(&model.UserLibrary{
		Model: model.Model{ID: uid},
	})
	switch verifyType {
	case "jwc":
		user.Update("jwc_verify", 0)
	case "my":
		user.Update("verify", 0)
	case "library":
		userLib.Update("verify", 0)
	}
	return err
}
