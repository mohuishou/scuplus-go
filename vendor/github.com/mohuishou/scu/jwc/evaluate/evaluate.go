package evaluate

import (
	"errors"
	"strings"
	"sort"

	"github.com/mohuishou/scu/jwc"
	"github.com/gocolly/colly"
	"github.com/mohuishou/scu/jwc/util"
	"github.com/PuerkitoBio/goquery"
	"log"
)

//Evaluate 评教信息列表
type Evaluate struct {
	EvaluateID   string  `json:"evaluate_id"`
	EvaluateType string  `json:"evaluate_type"`
	TeacherName  string  `json:"teacher_name"`
	TeacherID    string  `json:"teacher_id"`
	TeacherType  int     `json:"teacher_type"` // 0: 老师, 1:助教
	CourseName   string  `json:"course_name"`
	CourseID     string  `json:"course_id"`
	Status       int     `json:"status"`
	Star         float64 `json:"star"`    //平均分
	Comment      string  `json:"comment"` //评价内容
}

type EvaluateList []Evaluate

func (el EvaluateList) Len() int {
	return len(el)
}

func (el EvaluateList) Swap(i, j int) {
	el[i], el[j] = el[j], el[i]
}

func (el EvaluateList) Less(i, j int) bool {
	return el[i].Status < el[j].Status
}

func (e *Evaluate) getParams() map[string]string {
	return map[string]string{
		"wjbm":        e.EvaluateID,
		"wjmc":        util.Utf8ToGbk(e.EvaluateType),
		"bpr":         e.TeacherID,
		"bprm":        util.Utf8ToGbk(e.TeacherName),
		"pgnr":        e.CourseID,
		"pgnrm":       util.Utf8ToGbk(e.CourseName),
		"oper":        "wjShow",
		"pageSize":    "100",
		"page":        "1",
		"currentPage": "1",
	}
}

//EvaluateListURL 列表页面链接
const EvaluateListURL = jwc.DOMAIN + "/jxpgXsAction.do?oper=listWj&pageSize=100"

//EvaluateURL 评教页面链接
const EvaluateURL = jwc.DOMAIN + "/jxpgXsAction.do"

var nameToStar = map[string]float64{
	"非常同意":  5,
	"同意":    4,
	"基本同意":  3,
	"不同意":   2,
	"非常不同意": 1,
}

var StarToGradeTeacher = []string{"0.2", "0.6", "0.7", "0.8", "1"}
var StarToGradeZJ = []string{"0", "0.3", "0.6", "0.8", "1"}

// getEvaInfo 获取评教的详细信息
func (e *Evaluate) getEvaInfo(c *colly.Collector) (err error) {

	// 第一步：获取评教结果页面
	params := e.getParams()
	params["oper"] = "wjResultShow"

	starSum := 0.0
	count := 0.0
	c.OnHTML("#tblView > tbody > tr > td:nth-child(2) > table > tbody > tr> td > font", func(element *colly.HTMLElement) {
		if star, ok := nameToStar[strings.TrimSpace(element.Text)]; ok {
			starSum += star
			count += 1
		}
		if count == 0 {
			err = errors.New("获取分数信息失败！")
			return
		}
		e.Star = starSum / count
	})

	c.OnHTML("textarea[name=\"zgpj\"]", func(element *colly.HTMLElement) {
		e.Comment = strings.TrimSpace(element.Text)
	})

	c.Post(EvaluateURL, params)

	return err
}

// getStars 获取评教的平均分数
func (e *Evaluate) getStars(doc *goquery.Selection) error {
	starSum := 0.0
	count := 0.0
	doc.Find("tr td font").Each(func(i int, s *goquery.Selection) {
		res := strings.TrimSpace(s.Text())
		if star, ok := nameToStar[res]; ok {
			starSum += star
			count += 1
		}
	})
	log.Println(starSum, count)
	if count == 0 {
		return errors.New("获取分数信息失败！")
	}
	e.Star = starSum / count
	return nil
}

//AddEvaluate 评教
func AddEvaluate(c *colly.Collector, evaluate *Evaluate) (err error) {
	// 检查是否已经评教
	if evaluate.Status != 0 {
		return errors.New("您已评价：" + evaluate.CourseName + "-" + evaluate.TeacherName)
	}

	// 准备工作，必须先访问列表页
	c.Visit(EvaluateListURL)

	// 第一步：获取评教w问卷页面
	params := evaluate.getParams()
	var names []string
	c.OnHTML("#tblView > tbody > tr > td:nth-child(2) > table > tbody > tr > td > input[type=\"radio\"]:nth-child(1)", func(element *colly.HTMLElement) {
		name := element.Attr("name")
		if name != "" {
			names = append(names, name)
		}

		if len(names) == 0 {
			err = errors.New("问卷解析失败")
			return
		}
	})
	c.Post(EvaluateURL, params)
	c.Wait()
	if err != nil {
		return err
	}

	//发起评教请求
	params = map[string]string{}
	for _, name := range names {
		star := StarToGradeTeacher[int(evaluate.Star)-1]
		if evaluate.TeacherType == 1 {
			star = StarToGradeZJ[int(evaluate.Star)-1]
		}
		params[name] = "10_" + star
	}
	params["wjbm"] = evaluate.EvaluateID
	params["bpr"] = evaluate.TeacherID
	params["pgnr"] = evaluate.CourseID
	params["zgpj"] = util.Utf8ToGbk(evaluate.Comment)
	c.OnHTML("script", func(element *colly.HTMLElement) {
		if !strings.Contains(element.Text, "成功") {
			err = errors.New("评教失败！")
		}
	})
	c.Post(EvaluateURL+"?oper=wjpg", params)
	c.Wait()

	return err
}

// GetEvaList 获取评教数据
func GetEvaList(c *colly.Collector) (EvaluateList, error) {
	var evaluateList EvaluateList
	c.OnHTML("#user tbody td img", func(e *colly.HTMLElement) {
		eva := &Evaluate{}
		data := strings.Split(e.Attr("name"), "#@")
		if len(data) != 6 {
			return
		}
		eva.EvaluateID = data[0]
		eva.TeacherID = data[1]
		eva.TeacherName = data[2]
		eva.EvaluateType = data[3]
		eva.CourseName = data[4]
		eva.CourseID = data[5]
		eva.TeacherType = 1
		if strings.Contains(eva.EvaluateType, "学生") {
			eva.TeacherType = 0
		}
		switch strings.TrimSpace(e.Attr("title")) {
		case "评估":
			eva.Status = 0
		case "查看":
			eva.Status = 1
			eva.getEvaInfo(c)
		default:
			return
		}
		evaluateList = append(evaluateList, *eva)
	})
	c.Visit(EvaluateListURL)
	sort.Sort(evaluateList)
	return evaluateList, nil
}
