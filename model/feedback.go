package model

import (
	"log"
)

// Feedback 反馈
type Feedback struct {
	Model
	UserID uint   `json:"user_id"`
	Title  string `json:"title"`  // 反馈标题
	Number int    `json:"number"` // issue id
	Stat   string `json:"stat"`
	Tags   string `json:"tags"`
}

func UpdateFeedBack(num int, stat, tags string) error {
	fb := Feedback{}
	err := DB().Model(&fb).Where("number = ?", num).UpdateColumns(map[string]string{
		"stat": stat,
		"tags": tags,
	}).Error
	if err != nil {
		log.Println("update feedback err:", err)
	}
	return err
}
