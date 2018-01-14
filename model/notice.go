package model

// Notice 公告
type Notice struct {
	Model
	Title   string `json:"title"`
	Cover   string `json:"cover"` //封面链接
	Content string `gorm:"type:text" json:"content"`
	Status  int    `json:"status"` // 0: 正常, -1: 已关闭
}
