package model

import (
	"github.com/mohuishou/scu/ecard"
)

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
func UpdateEcard(uid uint) error {
	c, err := GetCollector(uid)
	if err != nil {
		return err
	}

	// 获取一卡通信息
	card, err := ecard.Get(c)

	// 更新一卡通余额
	if err := DB().Model(&UserInfo{UserID: uid}).Update("balance", card.Balance).Error; err != nil {
		return err
	}

	// 更新交易数据

	return nil
}
