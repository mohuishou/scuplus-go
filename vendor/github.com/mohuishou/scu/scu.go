package scu

import (
	"errors"
	"regexp"

	"github.com/gocolly/colly"
)

// NewCollector 新建一个采集器
func NewCollector(studentID, password string) (*colly.Collector, error) {
	c := colly.NewCollector()

	var logErr error
	c.OnHTML("script", func(e *colly.HTMLElement) {
		r, _ := regexp.Compile("错误")
		if r.MatchString(e.Text) {
			logErr = errors.New("[Error] 用户不存在或密码错误")
		}
	})

	err := c.Post("http://my.scu.edu.cn/userPasswordValidate.portal", map[string]string{
		"Login.Token1": studentID,
		"Login.Token2": password,
		"goto":         "http://my.scu.edu.cn/loginSuccess.portal",
		"gotoOnFail":   "http://my.scu.edu.cn/loginFailure.portal",
	})

	if err != nil {
		return nil, err
	}

	if logErr != nil {
		return nil, logErr
	}

	return c, nil
}
