package evaluate

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/gocolly/colly"
	"github.com/mohuishou/scu/jwc"
)

// ID id
type ID struct {
	EvaluateID string `json:"questionnaireCoding"`
	TeacherID  string `json:"evaluatedPeople"`
	CourseID   string `json:"evaluationContentNumber"`
}

// Info EvaluateInfo
type Info struct {
	EvaluateType string `json:"questionnaireName"`
}

//Evaluate 评教信息列表
type Evaluate struct {
	ID `json:"id"`
	// EvaluateID     string  `json:"evaluate_id"`
	// TeacherID      string  `json:"teacher_id"`
	// CourseID       string  `json:"course_id"`

	TeacherName string `json:"evaluatedPeople"`
	CourseName  string `json:"evaluationContent"`

	Info `json:"questionnaire"`

	TeacherType int     `json:"-"` // 0: 老师, 1:助教
	StatusStr   string  `json:"isEvaluated"`
	Status      int     `json:"-"`
	Star        float64 `json:"-"` //平均分
	Comment     string  `json:"-"` //评价内容
}

// 字段更新
func (e *Evaluate) update() {
	if e.StatusStr == "是" {
		e.Status = 1
	}

	if strings.Contains(e.EvaluateType, "助教") {
		e.TeacherType = 1
	}
}

// Evaluates Evaluates
type Evaluates []Evaluate

func (el Evaluates) Len() int {
	return len(el)
}

func (el Evaluates) Swap(i, j int) {
	el[i], el[j] = el[j], el[i]
}

func (el Evaluates) Less(i, j int) bool {
	return el[i].Status < el[j].Status
}

const baseURL = jwc.DOMAIN + "/student/teachingEvaluation/teachingEvaluation"

func (e *Evaluate) params() map[string]string {
	param := map[string]string{
		"evaluatedPeople":          e.TeacherName,
		"evaluatedPeopleNumber":    e.TeacherID,
		"questionnaireCode":        e.EvaluateID,
		"questionnaireName":        e.EvaluateType,
		"evaluationContentNumber":  e.CourseID,
		"evaluationContentContent": e.CourseID,
	}
	return param
}

// getInfo 获取评教的详细信息
func (e *Evaluate) getInfo(c *colly.Collector) (answers Answers, questions Questions, err error) {

	params := e.params()

	c.OnResponse(func(r *colly.Response) {
		// 只处理json
		contentType := r.Headers.Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			err = errors.New("页面类型错误！")
			return
		}

		type tmp struct {
			Answers  Answers `json:"selectedAnswerList"`
			PageDate struct {
				Questions Questions `json:"questionsList"`
			} `json:"pageDate"`
		}
		data := &tmp{}
		err = json.Unmarshal(r.Body, data)
		if err != nil {
			panic(err)
		}
		answers = data.Answers
		questions = data.PageDate.Questions
	})

	c.Post(baseURL+"/evaluationReusltPage", params)
	c.Wait()
	return answers, questions, err
}

func (e *Evaluate) getToken(c *colly.Collector) (token string) {

	c.OnHTML("input[name=\"tokenValue\"]", func(e *colly.HTMLElement) {
		token = e.Attr("value")
	})

	c.Post(jwc.DOMAIN+"/student/teachingEvaluation/teachingEvaluation/evaluationPage", e.params())
	c.Wait()
	return token
}

//AddEvaluate 评教
func AddEvaluate(c *colly.Collector, evaluate *Evaluate) (err error) {
	// 检查是否已经评教
	if evaluate.Status != 0 {
		return errors.New("您已评价：" + evaluate.CourseName + "-" + evaluate.TeacherName)
	}

	// 获取params
	_, questions, err := evaluate.getInfo(c)
	if err != nil {
		return err
	}
	params := questions.params(evaluate.Star, evaluate.TeacherType)

	// 获取token
	params["tokenValue"] = evaluate.getToken(c)
	if params["tokenValue"] == "" {
		return errors.New("token获取失败")
	}

	// 补充其余参数
	params["questionnaireCode"] = evaluate.EvaluateID
	params["evaluationContentNumber"] = evaluate.CourseID
	params["evaluatedPeopleNumber"] = evaluate.TeacherID
	params["zgpj"] = evaluate.Comment

	// 评教回调
	c.OnResponse(func(r *colly.Response) {
		if string(r.Body) != "success" {
			err = errors.New("评教失败")
		}
	})

	c.Post(jwc.DOMAIN+"/student/teachingEvaluation/teachingEvaluation/evaluation", params)
	c.Wait()

	return err
}

// GetEvaList 获取评教数据
func GetEvaList(c *colly.Collector) (evaluates Evaluates, err error) {
	c.OnResponse(func(r *colly.Response) {
		type tmp struct {
			Evaluates Evaluates `json:"data"`
		}
		data := &tmp{}
		err = json.Unmarshal(r.Body, data)
		evaluates = data.Evaluates
	})
	c.Visit(jwc.DOMAIN + "/student/teachingEvaluation/teachingEvaluation/search")
	c.Wait()

	if err != nil {
		return nil, err
	}

	for i, v := range evaluates {
		e := &v
		e.update()
		if e.Status == 1 {
			answer, _, _ := e.getInfo(c.Clone())
			e.Star = answer.average(e.TeacherType)
		}
		evaluates[i] = *e
	}

	return evaluates, nil
}
