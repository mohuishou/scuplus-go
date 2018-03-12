package model

import (
	"time"
)

// TermEvent 事件表，记录每个学期的当中的事件
type TermEvent struct {
	Model
	TermID    uint      `json:"term_id"`
	Name      string    `json:"name"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}
