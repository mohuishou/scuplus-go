package model

import (
	"log"
	"time"

	"github.com/mohuishou/scu/ecard"
)

// Ecard 一卡通交易数据
type Ecard struct {
	Model
	TransTime int64   `json:"trans_time"` // 交易时间
	Addr      string  `json:"addr"`       // 交易地点
	Money     float64 `json:"money"`      // 交易金额
	Balance   float64 `json:"balance"`    // 余额
}

func convertEcard(tran ecard.Transaction) Ecard {
	return Ecard{
		TransTime: tran.Time / 1000,
		Addr:      tran.Addr,
		Balance:   tran.Balance,
		Money:     tran.Money,
	}
}

// UpdateEcard 更新一卡通信息, 包括更新一卡通余额信息
// uid: 用户id
func UpdateEcard(uid uint) error {
	c, err := GetCollector(uid)
	if err != nil {
		return err
	}

	// 设置默认开始时间: 2个月内，结束时间: 当日
	end := time.Now()
	d, err := time.ParseDuration("60d")
	log.Println(err)
	start := end.Add(d)

	// 获取最后一条交易数据
	lastTrans := Ecard{}
	DB().Where("user_id = ?", uid).Order("timedesc").Last(&lastTrans)

	if lastTrans.ID != 0 {
		// 推后一天
		unix := (lastTrans.TransTime + (3600 * 24))
		start = time.Unix(unix, unix)
	}

	// 获取一卡通信息
	card, err := ecard.Get(c, start, end)

	// 更新一卡通余额
	if err := DB().Model(&UserInfo{UserID: uid}).Update("balance", card.Balance).Error; err != nil {
		return err
	}

	// 插入新的交易数据
	for _, v := range card.Transactions {
		eCard := convertEcard(v)
		if err := DB().Create(&eCard).Error; err != nil {
			log.Printf("[Error]: 更新一卡通数据错误,%s", err.Error())
			return err
		}
	}
	return nil
}
