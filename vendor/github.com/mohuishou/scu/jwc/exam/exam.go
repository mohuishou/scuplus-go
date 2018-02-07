package exam

import (
	"log"
	"reflect"
	"strings"

	"github.com/gocolly/colly"
	"github.com/mohuishou/scu/jwc"
)

// Exam 考试
type Exam struct {
	Name       string `json:"name"`
	Campus     string `json:"campus"`
	Building   string `json:"building"`
	Classroom  string `json:"classroom"`
	CourseName string `json:"course_name"`
	Week       string `json:"week"`
	Day        string `json:"day"`
	Date       string `json:"date"`
	Time       string `json:"time"`
	Site       string `json:"site"`
	Card       string `json:"card"`
	Comment    string `json:"comment"`
}

// Get 获取考试信息
func Get(c *colly.Collector) []Exam {
	exams := make([]Exam, 0)
	c.OnHTML("#user .odd", func(e *colly.HTMLElement) {
		exam := Exam{}
		v := reflect.ValueOf(&exam)
		elem := v.Elem()
		log.Println(elem.NumField())
		for k := 0; k < elem.NumField(); k++ {
			elem.Field(k).SetString(strings.TrimSpace(e.DOM.Find("td").Eq(k).Text()))
		}
		exams = append(exams, exam)
	})

	c.Visit(jwc.DOMAIN + "/ksApCxAction.do?oper=getKsapXx")
	return exams
}
