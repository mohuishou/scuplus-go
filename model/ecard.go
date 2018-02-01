package model

// Ecard 一卡通交易数据
type Ecard struct {
	Model
	Time    int64   `json:"time"`    // 交易时间
	Addr    string  `json:"addr"`    // 交易地点
	Money   float64 `json:"money"`   // 交易金额
	Balance float64 `json:"balance"` // 余额
}

// UpdateEcard 更新一卡通信息, 包括更新一卡通余额信息
// uid: 用户id
func UpdateEcard(uid uint) {

}
