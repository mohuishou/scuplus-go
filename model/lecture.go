package model

import "time"

type Lecture struct {
	Model
	Title     string    `json:"title" gorm:"unique_index:title_start"`      // 标题
	Time      string    `json:"time"`                                       // 时间
	StartTime time.Time `json:"start_time" gorm:"unique_index:title_start"` // 开始时间
	Address   string    `json:"address"`                                    // 地点
	Reporter  string    `json:"reporter"`                                   // 报告人
	College   string    `json:"college"`                                    // 举办学院
}
