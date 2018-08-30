package info

import (
	"reflect"

	"github.com/gocolly/colly"
	"github.com/mohuishou/scu/jwc"

	"strings"
)

//UserInfo 用户信息
type UserInfo struct {
	UID         string `json:"uid"`
	Name        string `json:"name"`
	NameEN      string `json:"name_en"`
	CardID      string `json:"card_id"`
	Sex         string `json:"sex"`
	StudentType string `json:"student_type"`
	Status      string `json:"status"`
	Nation      string `json:"nation"`
	Native      string `json:"native"`
	Birth       string `json:"birth"`
	Political   string `json:"political"`
	College     string `json:"college"`
	Major       string `json:"major"`
	Year        string `json:"year"`
	Class       string `json:"class"`
	Campus      string `json:"campus"`
}

// Get 获取用户信息
func Get(c *colly.Collector) (info UserInfo, err error) {
	userinfo := &UserInfo{}
	v := reflect.ValueOf(userinfo)
	elem := v.Elem()
	//对应关系
	eq := []int{0, 1, 3, 5, 6, 7, 9, 11, 12, 13, 14, 25, 26, 28, 29, 32}
	c.OnHTML("body", func(e *colly.HTMLElement) {
		for i := 0; i < elem.NumField(); i++ {
			s := e.DOM.Find("table td[width=\"275\"]").Eq(eq[i])
			elem.Field(i).SetString(strings.TrimSpace(s.Text()))
		}
	})
	c.Visit(jwc.DOMAIN + "/xjInfoAction.do?oper=xjxx")
	c.Wait()
	return *userinfo, nil
}
