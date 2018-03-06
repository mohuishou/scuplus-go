package ecard

import (
	"encoding/json"
	"errors"
	"sort"
	"time"

	"github.com/gocolly/colly"
)

// Transaction 交易记录
type Transaction struct {
	Time    int64   `json:"time"`    // 交易时间
	Addr    string  `json:"addr"`    // 交易地点
	Money   float64 `json:"money"`   // 交易金额
	Balance float64 `json:"balance"` // 余额
}

// Transactions 交易数据
type Transactions []Transaction

func (t Transactions) Len() int {
	return len(t)
}
func (t Transactions) Less(i, j int) bool {
	// 倒序排列
	return t[i].Time > t[j].Time
}
func (t Transactions) Swap(i, j int) {
	t[j], t[i] = t[i], t[j]
}

// Card 一卡通记录
type Card struct {
	Balance      float64      `json:"balance"`      // 余额
	Transactions Transactions `json:"transactions"` //交易记录
}

// Get 获取一卡通余额以及历史交易记录
func Get(c *colly.Collector, start, end time.Time) (Card, error) {
	c.Visit("http://ecard.scu.edu.cn/ajax/login/sso")

	trans := getTransactions(c, start, end)
	sort.Sort(trans)
	card := Card{
		Transactions: trans,
	}

	// balance, err := getBalance(c.Clone())

	// if err != nil {
	// 	return card, err
	// }

	// card.Balance = balance
	return card, nil
}

func getBalance(c *colly.Collector) (float64, error) {
	res := struct {
		JSONData struct {
			PageData []struct {
				Balance float64 `json:"balance"` // 余额
			} `json:"pageData"`
		} `json:"jsonData"`
	}{}
	c.OnResponse(func(r *colly.Response) {
		json.Unmarshal(r.Body, &res)
	})
	c.Post("http://ecard.scu.edu.cn/ajax/card/list.json", map[string]string{
		"int_start":   "0",
		"int_maxSize": "15",
	})
	if len(res.JSONData.PageData) < 1 {
		return 0, errors.New("获取余额失败")
	}
	return res.JSONData.PageData[0].Balance, nil

}

func getTransactions(c *colly.Collector, start, end time.Time) Transactions {
	transactions := make(Transactions, 0)
	res := struct {
		JSONData struct {
			PageData []struct {
				Time    int64   `json:"smtDealDateTime"` // 交易时间
				Addr    string  `json:"smtOrgName"`      // 交易地点
				Money   float64 `json:"smtTransMoney"`   // 交易金额
				Balance float64 `json:"smtOutMoney"`     // 余额
			} `json:"pageData"`
		} `json:"jsonData"`
	}{}
	c.OnResponse(func(r *colly.Response) {
		json.Unmarshal(r.Body, &res)
	})

	c.Post("http://ecard.scu.edu.cn/ajax/tran/list.json", map[string]string{
		"int_start":     "0",
		"int_maxSize":   "1000",
		"str_startTime": start.Format("2006-01-02"),
		"str_endTime":   end.Format("2006-01-02"),
	})

	for _, v := range res.JSONData.PageData {
		trans := Transaction{
			Time:    v.Time,
			Addr:    v.Addr,
			Money:   v.Money,
			Balance: v.Balance,
		}
		transactions = append(transactions, trans)
	}

	return transactions

}
