package schedule

import (
	"reflect"
	"strings"

	"github.com/mohuishou/scu/jwc/util"

	"github.com/PuerkitoBio/goquery"
	"github.com/mohuishou/scu/jwc"

	"github.com/gocolly/colly"
)

//Schedule 课程数据
type Schedule struct {
	Project    string   `json:"project"`
	CourseID   string   `json:"course_id"`
	CourseName string   `json:"course_name"`
	LessonID   string   `json:"lesson_id"`
	Credit     string   `json:"credit"`
	CourseType string   `json:"course_type"`
	ExamType   string   `json:"exam_type"`
	Teachers   []string `json:"teachers"`
	StudyWay   string   `json:"study_way"`
	ChooseType string   `json:"choose_type"`
	AllWeek    string   `json:"all_week"`
	Day        string   `json:"day"`
	Session    string   `json:"session"`
	Campus     string   `json:"campus"`
	Building   string   `json:"building"`
	Classroom  string   `json:"classroom"`
}

// Get 获取课程表
func Get(c *colly.Collector) (data []Schedule) {
	data = make([]Schedule, 0)

	//通过反射利用字段间的对应关系，来进行字段赋值
	c.OnHTML("body", func(e *colly.HTMLElement) {
		e.DOM.Find(".displayTag").Eq(1).Find("tr").Each(func(i int, sel *goquery.Selection) {
			schedule := &Schedule{}
			v := reflect.ValueOf(schedule)
			elem := v.Elem()
			td := sel.Find("td")
			index := 0
			k := 0
			t := elem.Type()

			//长度小于7说明，该课程为上一课程的不同时间段
			if td.Size() < 7 {
				k = 10
			}

			//获取数据
			for ; k < elem.NumField(); k++ {
				//跳过大纲日历
				if k == 8 {
					index++
				}

				s := td.Eq(index)

				switch t.Field(k).Name {
				case "Teachers":
					teachers := util.TeacherParse(s.Text())
					schedule.Teachers = teachers
				case "AllWeek":
					allWeek := util.WeekParse(s.Text())
					schedule.AllWeek = allWeek
				case "Session":
					session, _ := util.SessionParse(s.Text())
					schedule.Session = session
				default:
					elem.Field(k).SetString(strings.TrimSpace(s.Text()))
				}

				index++
			}

			//只有长度大于1，才说明这一行不是标题行
			if td.Size() > 1 {
				data = append(data, *schedule)
			}
		})
	})

	c.Visit(jwc.DOMAIN + "/xkAction.do?actionType=6")
	return data
}
