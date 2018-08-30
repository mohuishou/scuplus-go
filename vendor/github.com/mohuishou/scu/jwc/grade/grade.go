package grade

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/mohuishou/scu/jwc"

	"github.com/gocolly/colly"
)

//Grade 成绩
type Grade struct {
	// 新版新增，获取到之后会赋值给之前的字段
	ID struct {
		CourseID string `json:"courseNumber"`
		LessonID string `json:"coureSequenceNumber"`
	} `json:"id"`
	CourseID          string  `json:"-"`
	LessonID          string  `json:"-"`
	CourseName        string  `json:"courseName"`
	CourseEnglishName string  `json:"englishCourseName"`
	Credit            string  `json:"credit"`
	CourseType        string  `json:"courseAttributeName"`
	GradeShow         string  `json:"cj"`
	Grade             float64 `json:"courseScore"`
	GPA               float64 `json:"gradePointScore"`
	TermCode          string  `json:"termCode"` // 1: 秋季学期, 2: 春季学期
	Term              int     `json:"-"`        //0: 秋季学期, 1: 春季学期
	YearCode          string  `json:"academicYearCode"`
	Year              int     `json:"-"` //-1: 尚不及格，-2: 曾不及格
	TermName          string  `json:"-"`
}

// Grades 成绩列表
type Grades []Grade

// Term 一个学期的成绩
type Term struct {
	Grades    Grades  `json:"cjList"`
	TermName  string  `json:"cjlx"`
	AllCredit float64 `json:"yxxf"`
}

// Terms terms
type Terms []Term

func (terms Terms) getGrades() Grades {
	var grades Grades

	for _, term := range terms {
		grades = append(grades, term.getGrades()...)
	}

	return grades
}

func (term Term) getGrades() Grades {
	grades := make(Grades, len(term.Grades))
	for i, grade := range term.Grades {
		g := &grade
		g.update()
		g.TermName = term.TermName
		switch term.TermName {
		case "尚不及格":
			g.Year = -1
		case "曾不及格":
			g.Year = -2
		}
		grades[i] = *g
	}
	return grades
}

func (grade *Grade) update() {
	grade.CourseID = grade.ID.CourseID
	grade.LessonID = grade.ID.LessonID

	// 学期代码转换
	if grade.TermCode == "2" {
		grade.Term = 1
	}

	grade.Year, _ = strconv.Atoi(strings.Split(grade.YearCode, "-")[0])
}

var modes = map[string]string{
	"not_pass": "/student/integratedQuery/scoreQuery/unpassedScores/callback",
	"all":      "/student/integratedQuery/scoreQuery/allPassingScores/callback",
}

func get(c *colly.Collector, mode string) (Grades, error) {
	var (
		terms Terms
		err   error
	)
	c.OnResponse(func(r *colly.Response) {
		// 只处理json
		contentType := r.Headers.Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			return
		}

		type tmp struct {
			Terms Terms `json:"lnList"`
		}

		switch mode {
		case "all":
			data := &tmp{}
			err = json.Unmarshal(r.Body, data)
			terms = data.Terms
		case "not_pass":
			err = json.Unmarshal(r.Body, &terms)
		}
	})

	c.OnHTML(".alert", func(e *colly.HTMLElement) {
		err = errors.New(strings.TrimSpace(e.Text))
	})

	c.Visit(jwc.DOMAIN + modes[mode])
	c.Wait()
	if err != nil {
		return nil, err
	}

	return terms.getGrades(), nil
}

// GetNow 获取本学期成绩
func GetNow(c *colly.Collector) Grades {
	// TODO 暂时不知道api
	return nil
}

// GetALL 获取所有及格成绩
func GetALL(c *colly.Collector) (grades Grades, err error) {
	return get(c, "all")
}

// GetNotPass 获取所有不及格成绩
func GetNotPass(c *colly.Collector) (grades Grades, err error) {
	return get(c, "not_pass")
}
