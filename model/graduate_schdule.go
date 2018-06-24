package model

import (
	"errors"
	"fmt"

	"github.com/mohuishou/scu/ehall/schdule"
)

type GraduateSchedule struct {
	Model
	CourseID   string `json:"course_id"`   // 课程代码
	CourseName string `json:"course_name"` // 课程名
	Teacher    string `json:"teacher"`     // 教师名
	Address    string `json:"address"`     // 上课地点
	Day        int    `json:"day"`         // 上课星期
	SessionInt int    `json:"session_int"` // 上课节次
	AllWeek    string `json:"all_week"`    // 上课周次
	YearTerm   string `json:"year_term"`   // 学年学期
	Term       int    `json:"term"`        //0: 秋季学期, 1: 春季学期
	Year       int    `json:"year"`
	Session    string `json:"session"`
	UserID     uint   `json:"user_id"`
}

func UpdateGraduateSchedule(uid uint, year, term int) error {
	// 1. login in my.scu
	c, err := GetCollector(uid)
	if err != nil {
		return err
	}

	// 2. get new schedules
	schedules, err := schdule.Get(c, fmt.Sprintf("%d%d", year, term+1))
	if err != nil {
		return err
	}
	if len(schedules) == 0 {
		return errors.New("没有获取到新的数据，请查看教务处")
	}

	// 3. 删除所有该学期的课程
	DB().Unscoped().Delete(GraduateSchedule{}, GraduateSchedule{
		UserID: uid,
		Year:   year,
		Term:   term,
	})

	// 4. 保存
	for _, v := range schedules {
		s := &GraduateSchedule{}
		s.convert(v, uid)
		if err := DB().Create(s).Error; err != nil {
			return err
		}
	}
	return nil
}

// 类型转换
func (s *GraduateSchedule) convert(sch schdule.Schedule, uid uint) {
	s.CourseID = sch.CourseID
	s.CourseName = sch.CourseName
	s.Teacher = sch.Teacher
	s.Address = sch.Address
	s.Day = sch.Day
	s.SessionInt = sch.SessionInt
	s.AllWeek = sch.Week
	s.YearTerm = sch.YearTerm
	s.Term = sch.Term - 1
	s.Year = sch.Year
	s.Session = sch.Session
	s.UserID = uid
}
