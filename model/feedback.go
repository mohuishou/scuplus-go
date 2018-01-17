package model

// Feedback 反馈
type Feedback struct {
	Model
	UserID uint
	Title  string // 反馈标题
	Number int    // issue id
}
