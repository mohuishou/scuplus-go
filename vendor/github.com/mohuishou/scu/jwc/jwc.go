package jwc

import (
	"github.com/gocolly/colly"
	"github.com/mohuishou/scu"
)

// DOMAIN 教务处域名
const DOMAIN = "http://zhjw.scu.edu.cn"

// Login 登录教务处，获取已登录采集器
func Login(studentID, password string) (*colly.Collector, error) {
	c, err := scu.NewCollector(studentID, password)
	if err != nil {
		return nil, err
	}
	c.Visit(DOMAIN)
	return c, nil
}

// Logout 退出教务处
func Logout(c *colly.Collector) error {
	return c.Post(DOMAIN+"/logout.do", map[string]string{
		"loginType": "platformLogin",
	})
}
