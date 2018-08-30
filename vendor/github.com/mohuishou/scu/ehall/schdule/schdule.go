package schdule

import (
	"strconv"

	"github.com/gocolly/colly"
	"github.com/json-iterator/go"
	"github.com/mohuishou/scu/ehall"
	"github.com/mohuishou/scu/jwc/util"
)

//Schedule 成绩
type Schedule struct {
	CourseID   string `json:"KCDM"`   // 课程代码
	CourseName string `json:"KCMC"`   // 课程名
	Teacher    string `json:"JSXM"`   // 教师名
	Address    string `json:"JASMC"`  // 上课地点
	Day        int    `json:"XQ"`     // 上课星期
	SessionInt int    `json:"KSJCDM"` // 上课节次
	Week       string `json:"ZCMC"`   // 上课周次
	YearTerm   string `json:"XNXQDM"` // 学年学期
	Term       int    `json:"term"`   //0: 秋季学期, 1: 春季学期
	Year       int    `json:"year"`
	Session    string `json:"session"`
}

func Get(c *colly.Collector, yearTerm string) (schedules []Schedule, err error) {
	defer ehall.Logout(c)
	// 需要先访问前置界面
	c.Visit(ehall.DOMAIN + "/appShow?appId=4979568947762216")

	// 获取成绩
	c.OnResponse(func(response *colly.Response) {
		data := struct {
			Schedules []Schedule `json:"xspkjgList"`
		}{}
		err = jsoniter.Unmarshal(response.Body, &data)
		schedules = data.Schedules
	})

	err = c.Post(ehall.DOMAIN+"/gsapp/sys/wdkbapp/wdkcb/xspkjgQuery.do", map[string]string{
		"XNXQDM": yearTerm,
	})
	c.Wait()

	if err != nil {
		return nil, err
	}

	// 数据处理
	tmp := map[string]Schedule{}
	for _, v := range schedules {
		// 年份学期处理
		v.convertYearTerm()

		// 周次处理
		v.Week = util.WeekParse(v.Week)

		// 节次处理
		v.Session = strconv.Itoa(v.SessionInt)
		s, ok := tmp[v.CourseID]
		if !ok {
			tmp[v.CourseID] = v
		} else {
			if v.SessionInt > s.SessionInt {
				s.Session = s.Session + "," + v.Session
			} else {
				s.Session = v.Session + "," + s.Session
			}
			tmp[v.CourseID] = s
		}
	}

	schedules = []Schedule{}
	for _, v := range tmp {
		schedules = append(schedules, v)
	}

	return schedules, nil
}

func (s *Schedule) convertYearTerm() {
	t, _ := strconv.Atoi(s.YearTerm)
	s.Year = t / 10
	s.Term = t % 10
}
