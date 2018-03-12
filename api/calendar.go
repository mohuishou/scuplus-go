package api

import (
	"time"

	"github.com/mohuishou/scuplus-go/model"

	"github.com/kataras/iris"
)

// GetTerm 获取当前学期
func GetTerm(ctx iris.Context) {
	now := time.Now()
	var term model.Term
	err := model.DB().Where("start_time < ? and end_time > ?", now, now).Find(&term).Error
	if err != nil {
		Error(ctx, 80001, "获取本学期错误", err)
		return
	}
	Success(ctx, "获取成功", term)
}

// GetTermEvents 获取当前学期的所有事件
func GetTermEvents(ctx iris.Context) {
	now := time.Now()
	var term model.Term
	err := model.DB().Where("start_time < ? and end_time > ?", now, now).Find(&term).Error
	if err != nil {
		Error(ctx, 80001, "获取本学期错误", err)
		return
	}
	events := []model.TermEvent{}
	model.DB().Where("term_id = ?", term.ID).Find(&events)
	Success(ctx, "获取成功", map[string]interface{}{
		"term":   term,
		"events": events,
	})
}
