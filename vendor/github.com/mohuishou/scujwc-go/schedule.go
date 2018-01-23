package scujwc

import (
	"errors"
	"regexp"
	"strings"

	"reflect"

	"strconv"

	"github.com/PuerkitoBio/goquery"
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

//Schedule 课程表
func (j *Jwc) Schedule() (data []Schedule, err error) {
	data, err = getSchedule(*j)
	return data, err
}

func getSchedule(j Jwc) (data []Schedule, err error) {
	data = make([]Schedule, 0)
	url := DOMAIN + "/xkAction.do"
	doc, err := j.jPost(url, "actionType=6")
	if err != nil {
		return nil, err
	}

	//通过反射利用字段间的对应关系，来进行字段赋值
	schedule := &Schedule{}
	v := reflect.ValueOf(schedule)
	elem := v.Elem()

	doc.Find(".displayTag").Eq(1).Find("tr").Each(func(i int, sel *goquery.Selection) {
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
				teachers := teacherParse(s.Text())
				schedule.Teachers = teachers
			case "AllWeek":
				allWeek := weekParse(s.Text())
				schedule.AllWeek = allWeek
			case "Session":
				session, _ := sessionParse(s.Text())
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
	return data, nil
}

//教师解析，返回包含每个教师名字的数组
func teacherParse(t string) (teachers []string) {
	t = strings.TrimSpace(t)
	teachers = strings.Split(t, " ")
	return teachers
}

//上课时间解析
func weekParse(w string) (allWeek string) {
	re, _ := regexp.Compile(`[1-9]\d*|单|双`)
	s := re.FindAllStringSubmatch(w, -1)
	if len(s) == 1 {
		if s[0][0] == "单" {
			return "1,3,5,7,9,11,13,15,17"
		} else if s[0][0] == "双" {
			return "2,4,6,8,10,12,14,16,18"
		}
	} else if len(s) == 2 {
		start, _ := strconv.Atoi(s[0][0])
		end, _ := strconv.Atoi(s[1][0])
		for i := start; i < end; i++ {
			is := strconv.Itoa(i)
			allWeek = allWeek + is + ","
		}
		allWeek = allWeek + s[1][0]
	}
	return allWeek
}

func sessionParse(session string) (data string, err error) {
	session = strings.TrimSpace(session)
	sessions := strings.Split(session, "~")
	if len(sessions) != 2 {
		//todo:解析
		return "", errors.New("错误")
	}
	start, _ := strconv.Atoi(sessions[0])
	end, _ := strconv.Atoi(sessions[1])
	for i := start; i < end; i++ {
		s := strconv.Itoa(i)
		data = data + s + ","
	}
	data = data + sessions[1]
	return data, nil
}
