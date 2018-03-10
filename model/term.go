package model

import (
	"time"
)

// Term 学期表，记录每个学期的开学时间与放假时间
type Term struct {
	Model
	Name      string    `json:"name"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}
