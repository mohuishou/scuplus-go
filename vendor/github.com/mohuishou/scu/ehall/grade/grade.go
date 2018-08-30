package grade

import (
	"strconv"

	"github.com/gocolly/colly"
	"github.com/json-iterator/go"
	"github.com/mohuishou/scu/ehall"
)

const (
	TermAutumn = 0
	TermSpring = 1
)

//Grade 成绩
type Grade struct {
	CourseID   string  `json:"KCDM"`    // 课程代码
	CourseName string  `json:"KCMC"`    // 课程名
	Credit     float64 `json:"XF"`      // 学分
	CourseType string  `json:"KCLBMC"`  // 课程类别
	GradeShow  string  `json:"CJXSZ"`   // 成绩显示值
	GPA        float64 `json:"JDZ"`     // 绩点
	Grade      float64 `json:"DYBFZCJ"` // 百分制成绩
	ExamType   string  `json:"KSXZDM"`  // 考试性质
	YearTerm   string  `json:"XNXQDM"`  // 学年学期
	Term       int     `json:"term"`    //0: 秋季学期, 1: 春季学期
	Year       int     `json:"year"`
}

func Get(c *colly.Collector) (grades []Grade, err error) {
	defer ehall.Logout(c)

	// 需要先访问前置界面
	c.Visit(ehall.DOMAIN + "/appShow?appId=5094115980385668")

	// 获取成绩
	c.OnResponse(func(response *colly.Response) {
		data := struct {
			Datas struct {
				Xscicx struct {
					Grades []Grade `json:"rows"`
				} `json:"xscjcx"`
			} `json:"datas"`
		}{}
		err = jsoniter.Unmarshal(response.Body, &data)
		grades = data.Datas.Xscicx.Grades
	})

	err = c.Post(ehall.DOMAIN+"/gsapp/sys/wdcjapp/modules/wdcj/xscjcx.do", map[string]string{
		"pageSize":   "500",
		"pageNumber": "1",
	})
	c.Wait()

	if err != nil {
		return nil, err
	}

	// 处理成绩
	for k, g := range grades {
		t, _ := strconv.Atoi(g.YearTerm)
		grades[k].Year = t / 10
		grades[k].Term = t % 10
	}

	return grades, nil
}
