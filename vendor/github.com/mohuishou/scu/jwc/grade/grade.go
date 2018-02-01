package grade

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/mohuishou/scu/jwc"

	"github.com/gocolly/colly"

	"github.com/PuerkitoBio/goquery"
)

//Grade 成绩
type Grade struct {
	CourseID          string `json:"course_id"`
	LessonID          string `json:"lesson_id"`
	CourseName        string `json:"course_name"`
	CourseEnglishName string `json:"course_english_name"`
	Credit            string `json:"credit"`
	CourseType        string `json:"course_type"`
	Grade             string `json:"grade"`
	Term              int    `json:"term"` //0: 秋季学期, 1: 春季学期
	Year              int    `json:"year"`
	TermName          string `json:"term_name"`
}

// Grades 成绩列表
type Grades []Grade

func get(doc *goquery.Selection, year, term int, termName string) Grades {
	grades := make(Grades, 0)
	//抓取数据
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return
		}
		grade := Grade{Term: term, Year: year, TermName: termName}
		v := reflect.ValueOf(&grade)
		elem := v.Elem()
		for k := 0; k < elem.NumField(); k++ {
			if k > 6 {
				break
			}
			elem.Field(k).SetString(strings.TrimSpace(s.Find("td").Eq(k).Text()))
		}
		grades = append(grades, grade)
	})
	return grades
}

// GetNow 获取本学期成绩
func GetNow(c *colly.Collector) Grades {
	var grades Grades
	c.OnHTML("#user", func(e *colly.HTMLElement) {
		grades = get(e.DOM, 0, 0, "本学期成绩")
	})
	c.Visit(jwc.DOMAIN + "/bxqcjcxAction.do?pageSize=200")
	return grades
}

// GetALL 获取所有及格成绩
func GetALL(c *colly.Collector) Grades {
	var grades Grades
	tmps := map[string][]string{
		"term":  []string{},
		"year":  []string{},
		"title": []string{},
	}
	r, _ := regexp.Compile(`(\d+)-\d+学年(.)`)
	c.OnHTML("table b", func(e *colly.HTMLElement) {
		tmps["title"] = append(tmps["title"], e.Text)
		res := r.FindAllStringSubmatch(e.Text, -1)
		if len(res[0]) == 3 {
			tmps["year"] = append(tmps["year"], res[0][1])
			tmps["term"] = append(tmps["term"], res[0][2])
		}
	})
	i := 0
	c.OnHTML("#user", func(e *colly.HTMLElement) {
		year, term := 0, 0
		title := ""
		if len(tmps["year"]) > i {
			title = tmps["title"][i]
			year, _ = strconv.Atoi(tmps["year"][i])
			if tmps["term"][i] == "春" {
				term = 1
			}
			i++
		}
		grades = append(grades, get(e.DOM, year, term, title)...)
	})
	c.Visit(jwc.DOMAIN + "/gradeLnAllAction.do?type=ln&oper=qbinfo&lnxndm")
	return grades
}

// GetNotPass 获取所有不及格成绩
func GetNotPass(c *colly.Collector) Grades {
	termNames := []string{"尚不及格", "曾不及格"}
	i := 0
	var grades Grades
	c.OnHTML("#user", func(e *colly.HTMLElement) {
		grades = append(grades, get(e.DOM, 0, 0, termNames[i])...)
		i++
	})
	c.Visit(jwc.DOMAIN + "/gradeLnAllAction.do?type=ln&oper=bjg")
	return grades
}
