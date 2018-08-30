package ehall

import (
	"github.com/gocolly/colly"
	"github.com/mohuishou/scu"
)

const DOMAIN = "http://ehall.scu.edu.cn"

func Login(studentID, password string) (*colly.Collector, error) {
	c, err := scu.NewCollector(studentID, password)
	if err != nil {
		return nil, err
	}

	c.Visit("http://ehall.scu.edu.cn/login?service=http://ehall.scu.edu.cn/new/index.html")
	return c, nil
}

func Logout(c *colly.Collector) {
	c.Visit("http://ehall.scu.edu.cn/logout?service=http://ehall.scu.edu.cn/new/index.html")
}
