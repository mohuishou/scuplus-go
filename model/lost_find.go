package model

// LostFind 失物招领
type LostFind struct {
	Model
	UserID   uint   `json:"user_id"`  // 用户id
	Title    string `json:"title"`    // 标题
	Pictures string `json:"pictures"` // 截图链接
	Info     string `json:"info"`     // 信息
	Address  string `json:"address"`  // 地点
	Contact  string `json:"contact"`  // 联系方式
	Category string `json:"category"` // 分类: 一卡通,其他,遗失
}
