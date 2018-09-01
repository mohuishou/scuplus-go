package schedule

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mohuishou/scu/jwc"

	"github.com/gocolly/colly"
)

//Schedule 课程数据
type Schedule struct {
	// 新版新增，获取到之后会赋值给之前的字段
	ID struct {
		CourseID string `json:"courseNumber"`
		LessonID string `json:"coureSequenceNumber"`
	} `json:"id"`
	Project      string        `json:"programPlanName"`
	CourseID     string        `json:"-"`
	CourseName   string        `json:"courseName"`
	LessonID     string        `json:"-"`
	Credit       float64       `json:"unit"`
	CourseType   string        `json:"coursePropertiesName"`
	ExamType     string        `json:"examTypeName"`
	TeacherStr   string        `json:"attendClassTeacher"`
	Teachers     []string      `json:"-"`
	StudyWay     string        `json:"courseCategoryName"`
	ChooseType   string        `json:"selectCourseStatusName"`
	TimeAndAddrs []TimeAndAddr `json:"timeAndPlaceList"`
	AllWeek      string        `json:"-"`
	Day          int           `json:"-"`
	Session      string        `json:"-"`
	Campus       string        `json:"-"`
	Building     string        `json:"-"`
	Classroom    string        `json:"-"`
}

func (s Schedule) update() []Schedule {
	s.Teachers = strings.Split(s.TeacherStr, " ")
	s.CourseID = s.ID.CourseID
	s.LessonID = s.ID.LessonID

	schedules := make([]Schedule, len(s.TimeAndAddrs))
	for i, v := range s.TimeAndAddrs {
		s.AllWeek = getAllWeek(v.AllWeek)
		s.Session = getSession(v.StartSession, v.ContinueSession)
		s.Day = v.Day
		s.Campus = v.Campus
		s.Building = v.Building
		s.Classroom = v.Classroom
		schedules[i] = s
	}
	return schedules
}

func getAllWeek(week string) string {
	weeks := strings.Split(week, "")
	allWeek := ""
	for i, v := range weeks {
		if v == "1" {
			allWeek = fmt.Sprintf("%s%d,", allWeek, i+1)
		}
	}
	return strings.Trim(allWeek, ",")
}

func getSession(start, n int) (session string) {
	for i := 0; i < n; i++ {
		session = fmt.Sprintf("%s%d,", session, start+i)
	}
	return strings.Trim(session, ",")
}

// TimeAndAddr 上课时间地点
type TimeAndAddr struct {
	AllWeek         string `json:"classWeek"`
	Day             int    `json:"classDay"`
	StartSession    int    `json:"classSessions"`
	ContinueSession int    `json:"continuingSession"`
	Campus          string `json:"campusName"`
	Building        string `json:"teachingBuildingName"`
	Classroom       string `json:"classroomName"`
}

// Get 获取课程表
func Get(c *colly.Collector) (schedules []Schedule, err error) {
	c.OnResponse(func(r *colly.Response) {
		type tmp struct {
			DataList []struct {
				Schedules []Schedule `json:"selectCourseList"`
			} `json:"dateList"`
		}

		data := &tmp{}

		err = json.Unmarshal(r.Body, data)

		if len(data.DataList) == 1 {
			schedules = data.DataList[0].Schedules
		}
	})

	c.Visit(jwc.DOMAIN + "/student/courseSelect/thisSemesterCurriculum/ajaxStudentSchedule/callback")
	c.Wait()

	if err != nil {
		return nil, err
	}

	var newSch []Schedule
	for _, schedule := range schedules {
		newSch = append(newSch, schedule.update()...)
	}
	return newSch, nil
}
