package course

import (
	"log"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/mohuishou/scu/jwc"

	"github.com/mohuishou/scu/jwc/util"

	"github.com/gocolly/colly"

	"github.com/PuerkitoBio/goquery"
)

// Course 课程获取
type Course struct {
	College     string   // 学院
	CourseID    string   // 课程号
	Name        string   // 课程名
	LessonID    string   // 课序号
	Credit      float64  // 学分
	ExamType    string   // 考试类型
	Teachers    []string //教师
	AllWeek     string   // 周次: 1,2,3,4
	Day         int      // 星期
	Session     string   // 节次 1,2
	Campus      string   // 校区
	Building    string   // 教学楼
	Classroom   string   // 教室
	Max         int      // 课容量
	StudentNo   int      // 学生数
	CourseLimit string   // 选课限制说明
}

func getParamsStr(p url.Values) string {
	params, _ := url.ParseQuery("kch=&kcm=&jsm=&xsjc=&skxq=&skjc=&xaqh=&jxlh=&jash=&pageSize=20&showColumn=kkxsjc%23%BF%AA%BF%CE%CF%B5&showColumn=kch%23%BF%CE%B3%CC%BA%C5&showColumn=kcm%23%BF%CE%B3%CC%C3%FB&showColumn=kxh%23%BF%CE%D0%F2%BA%C5&showColumn=xf%23%D1%A7%B7%D6&showColumn=kslxmc%23%BF%BC%CA%D4%C0%E0%D0%CD&showColumn=skjs%23%BD%CC%CA%A6&showColumn=zcsm%23%D6%DC%B4%CE&showColumn=skxq%23%D0%C7%C6%DA&showColumn=skjc%23%BD%DA%B4%CE&showColumn=xqm%23%D0%A3%C7%F8&showColumn=jxlm%23%BD%CC%D1%A7%C2%A5&showColumn=jasm%23%BD%CC%CA%D2&showColumn=bkskrl%23%BF%CE%C8%DD%C1%BF&showColumn=xss%23%D1%A7%C9%FA%CA%FD&showColumn=xkxzsm%23%D1%A1%BF%CE%CF%DE%D6%C6%CB%B5%C3%F7&pageNumber=0&actionType=1")
	for k := range p {
		params.Set(k, p.Get(k))
	}
	return params.Encode()
}

// Get 获取本学期课程信息
func Get(c *colly.Collector, params url.Values) []Course {

	courseList := make([]Course, 0)

	c.OnHTML("#user > tbody", func(e *colly.HTMLElement) {
		e.DOM.Find("tr").Each(func(i int, s *goquery.Selection) {
			c := Course{}
			v := reflect.ValueOf(&c)
			elem := v.Elem()
			typeOfCource := elem.Type()
			for k := 0; k < elem.NumField(); k++ {
				switch typeOfCource.Field(k).Name {
				case "Credit":
					credit, err := strconv.ParseFloat(strings.TrimSpace(s.Find("td").Eq(k).Text()), 64)
					if err != nil {
						log.Println("[Error] 学分获取失败")
						return
					}
					elem.Field(k).SetFloat(credit)
				case "Day", "Max", "StudentNo":
					val, err := strconv.Atoi(strings.TrimSpace(s.Find("td").Eq(k).Text()))
					if err != nil {
						log.Println("[Error] 字符串转int失败")
						val = 0
					}
					elem.Field(k).SetInt(int64(val))
				case "Teachers":
					c.Teachers = util.TeacherParse(strings.TrimSpace(s.Find("td").Eq(k).Text()))
				case "AllWeek":
					allWeek := util.WeekParse(strings.TrimSpace(s.Find("td").Eq(k).Text()))
					elem.Field(k).SetString(allWeek)
				case "Session":
					session, _ := util.SessionParse(strings.TrimSpace(s.Find("td").Eq(k).Text()))
					elem.Field(k).SetString(session)
				default:
					elem.Field(k).SetString(strings.TrimSpace(s.Find("td").Eq(k).Text()))
				}
			}
			courseList = append(courseList, c)
		})
	})

	c.Visit(jwc.DOMAIN + "/courseSearchAction.do?" + getParamsStr(params))
	return courseList
}
