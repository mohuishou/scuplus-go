package jwc

import (
	"errors"

	"github.com/gocolly/colly"
)

// DOMAIN 教务处域名
const DOMAIN = "http://zhjw.scu.edu.cn"

// Login 登录教务处，获取已登录采集器
func Login(studentID, password string) (*colly.Collector, error) {
	c := colly.NewCollector()
	// 判定是否登录失败
	var logErr error
	c.OnHTML("font[color=\"#990000\"]", func(e *colly.HTMLElement) {
		if e.Text != "" {
			logErr = errors.New(e.Text)
		}
	})

	// 尝试登录
	if err := c.Post(DOMAIN+"/loginAction.do", map[string]string{
		"zjh": studentID,
		"mm":  password,
	}); err != nil {
		return nil, err
	}

	// 如果登录失败，返回登录信息
	if logErr != nil {
		return nil, logErr
	}

	return c, nil
}

// Logout 退出教务处
func Logout(c *colly.Collector) error {
	return c.Post(DOMAIN+"/logout.do", map[string]string{
		"loginType": "platformLogin",
	})
}
