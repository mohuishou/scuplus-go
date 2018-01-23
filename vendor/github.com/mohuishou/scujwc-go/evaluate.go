package scujwc

import (
	"errors"
	"log"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

func (e *Evaluate) getParams() url.Values {
	params := url.Values{}
	params.Set("wjbm", e.EvaluateID)
	params.Set("wjmc", Utf8ToGbk(e.EvaluateType))
	params.Set("bpr", e.TeacherID)
	params.Set("bprm", Utf8ToGbk(e.TeacherName))
	params.Set("pgnr", e.CourseID)
	params.Set("pgnrm", Utf8ToGbk(e.CourseName))
	params.Set("oper", "wjShow")
	params.Set("pageSize", "20")
	params.Set("page", "1")
	params.Set("currentPage", "1")
	return params
}

//EvaluateListURL 列表页面链接
const EvaluateListURL = DOMAIN + "/jxpgXsAction.do?oper=listWj&pageSize=100"

//EvaluateURL 评教页面链接
const EvaluateURL = DOMAIN + "/jxpgXsAction.do"

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
func (e *Evaluate) getEvaInfo(j *Jwc) error {

	// 第一步：获取评教结果页面
	params := e.getParams()
	params.Set("oper", "wjResultShow")
	doc, err := j.post(EvaluateURL, params.Encode())
	if err != nil {
		log.Println(err)
		return err
	}

	// 第二步获取评教的平均分数
	err = e.getStars(doc)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	//第三步获取评教内容
	s := doc.Find("textarea[name=\"zgpj\"]")
	e.Comment = strings.TrimSpace(s.Text())
	return nil
}

// getStars 获取评教的平均分数
func (e *Evaluate) getStars(doc *goquery.Document) error {
	starSum := 0.0
	count := 0.0
	doc.Find("#tblView table tbody tr td font").Each(func(i int, s *goquery.Selection) {
		res := strings.TrimSpace(s.Text())
		if star, ok := nameToStar[res]; ok {
			starSum += star
			count += 1
		}
	})
	if count == 0 {
		return errors.New("获取分数信息失败！")
	}
	e.Star = starSum / count
	return nil
}

// getStarName 获取所有分数项radio的name
func getStarName(doc *goquery.Document) ([]string, error) {
	names := make([]string, 0)
	doc.Find("#tblView table tbody td input[type=\"radio\"]").Each(func(i int, s *goquery.Selection) {
		name, exist := s.Attr("name")
		if !exist {
			return
		}
		for _, v := range names {
			if v == name {
				return
			}
		}
		names = append(names, name)
	})
	return names, nil
}

// Evaluate 评教
func (j *Jwc) Evaluate(evaluate *Evaluate) error {
	defer j.Logout()

	// 检查是否已经评教
	if evaluate.Status != 0 {
		return errors.New("您已评价：" + evaluate.CourseName + "-" + evaluate.TeacherName)
	}

	// 第一步：获取评教w问卷页面
	params := evaluate.getParams()
	doc, err := j.post(EvaluateURL, params.Encode())
	if err != nil {
		log.Println(err)
		return err
	}

	//获取需要评教的星级
	names, err := getStarName(doc)
	if err != nil || len(names) == 0 {
		log.Println(err)
		err = errors.New(err.Error() + "star names 获取失败")
		return err
	}

	//构造评教参数
	params = url.Values{}
	log.Println("names:", names)
	for _, name := range names {
		star := StarToGradeTeacher[int(evaluate.Star)-1]
		if evaluate.TeacherType == 1 {
			star = StarToGradeZJ[int(evaluate.Star)-1]
		}
		params.Set(name, "10_"+star)
	}
	params.Set("wjbm", evaluate.EvaluateID)
	params.Set("bpr", evaluate.TeacherID)
	params.Set("pgnr", evaluate.CourseID)
	params.Set("zgpj", Utf8ToGbk(evaluate.Comment))

	log.Println(params.Encode())

	doc, err = j.post(EvaluateURL+"?oper=wjpg", params.Encode())
	if err != nil {
		log.Println(err)
		return err
	}
	if !strings.Contains(doc.Find("script").Text(), "成功") {
		return errors.New("评教失败！")
	}

	return nil
}

// GetEvaList 获取评教数据
func (j *Jwc) GetEvaList() ([]Evaluate, error) {
	defer j.Logout()
	doc, err := j.get(EvaluateListURL, "")
	evaluateList := make([]Evaluate, 0)
	if err != nil {
		return nil, err
	}
	doc.Find("#user tbody td img").Each(func(i int, selection *goquery.Selection) {
		eva := &Evaluate{}
		val, exist := selection.Attr("name")
		data := strings.Split(val, "#@")
		if !exist || len(data) != 6 {
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

		val, exist = selection.Attr("title")
		if !exist {
			return
		}
		switch strings.TrimSpace(val) {
		case "评估":
			eva.Status = 0
		case "查看":
			eva.Status = 1
			eva.getEvaInfo(j)
		default:
			return
		}
		evaluateList = append(evaluateList, *eva)
	})
	return evaluateList, nil
}
