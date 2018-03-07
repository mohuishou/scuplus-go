package model

// Feedback 反馈
type Feedback struct {
	Model
	UserID uint   `json:"user_id"`
	Title  string `json:"title"`  // 反馈标题
	Number int    `json:"number"` // issue id
}
