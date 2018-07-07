package model

// Notice 公告
type Notice struct {
	Model
	Title    string `json:"title"`
	Cover    string `json:"cover"`    //封面链接
	Abstract string `json:"abstract"` // 摘要
	Content  string `gorm:"type:text" json:"content"`
	Status   int    `json:"status"`                  // 0: 正常, -1: 已关闭
	Newest   int    `json:"newest" gorm:"default:0"` // 1: 最新通知，弹窗提醒
}
