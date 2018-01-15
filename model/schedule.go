package model

import (
	"errors"
	"strings"

	"github.com/mohuishou/scujwc-go"
)

// Schedule 课程表
type Schedule struct {
	Model
	Project      string `json:"project"`
	CourseID     string `json:"course_id"`
	CourseName   string `json:"course_name"`
	LessonID     string `json:"lesson_id"`
	Credit       string `json:"credit"`
	CourseType   string `json:"course_type"`
	ExamType     string `json:"exam_type"`
	Teachers     string `json:"teachers"`
	StudyWay     string `json:"study_way"`
	ChooseStatus string `json:"choose_status"`
	AllWeek      string `json:"all_week"`
	Day          string `json:"day"`
	Session      string `json:"session"`
	Campus       string `json:"campus"`
	Building     string `json:"building"`
	Classroom    string `json:"classroom"`
	Term         string // 学期，例如：2017-2018年春季学期
	UserID       uint
}

// ScheduleList 课程表数组
type ScheduleList []Schedule

// GetSchedules 获取某个用户的某个学期课程表信息
// TODO: 待完善
func GetSchedules(userID uint, term string) []Schedule {
	schedules := []Schedule{}
	DB().Where(Schedule{UserID: userID, Term: term}).Find(&schedules)
	return schedules
}

// UpdateSchedules 更新课程表
func UpdateSchedules(userID uint, term string) error {
	jwc, err := GetJwc(userID)
	if err != nil {
		return err
	}

	schedules, err := jwc.Schedule()
	if err != nil {
		return err
	}

	if len(schedules) < 1 {
		return errors.New("没有获取到新的数据，请查看教务处")
	}

	// 删除所有的数据，软删除只保留一个版本
	// TODO: 待后期优化
	DB().Unscoped().Where("deleted_at IS NOT NULL").Delete(Schedule{}, Schedule{UserID: userID, Term: term})

	if err := DB().Delete(Schedule{}, Schedule{UserID: userID, Term: term}).Error; err != nil {
		return err
	}

	for _, schedule := range schedules {
		newSchedule := Schedule{
			Term:         term,
			UserID:       userID,
			Project:      schedule.Project,
			CourseID:     schedule.CourseID,
			CourseName:   schedule.CourseName,
			LessonID:     schedule.LessonID,
			Credit:       schedule.Credit,
			CourseType:   schedule.CourseType,
			ExamType:     schedule.ExamType,
			Teachers:     strings.Join(schedule.Teachers, ","),
			StudyWay:     schedule.StudyWay,
			ChooseStatus: schedule.ChooseType,
			AllWeek:      schedule.AllWeek,
			Day:          schedule.Day,
			Session:      schedule.Session,
			Campus:       schedule.Campus,
			Building:     schedule.Building,
			Classroom:    schedule.Classroom,
		}

		if err := DB().Create(&newSchedule).Error; err != nil {
			return err
		}
	}
	return nil
}

// deleteMore 删除多余的数据
// 传入最新抓取的数据
func (sl ScheduleList) deleteMore(s []scujwc.Schedule) error {
	deleteID := []uint{}
OLD:
	for _, schedule := range sl {
		for _, sch := range s {
			if schedule.CourseID == sch.CourseID && schedule.LessonID == sch.LessonID && schedule.Day == sch.Day && schedule.Session == sch.Session {
				continue OLD
			}
		}

		// 最新数据没有找到课程
		deleteID = append(deleteID, schedule.ID)
	}

	if err := DB().Unscoped().Delete(Schedule{}, "id in (?)", deleteID).Error; err != nil {
		return err
	}
	return nil
}
