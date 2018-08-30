package model

import (
	"strings"

	"github.com/mohuishou/scu/jwc"
	"github.com/mohuishou/scu/jwc/schedule"
)

// Schedule 课程表
type Schedule struct {
	Model
	Project      string  `json:"project"`
	CourseID     string  `json:"course_id"`
	CourseName   string  `json:"course_name"`
	LessonID     string  `json:"lesson_id"`
	Credit       float64 `json:"credit"`
	CourseType   string  `json:"course_type"`
	ExamType     string  `json:"exam_type"`
	Teachers     string  `json:"teachers"`
	StudyWay     string  `json:"study_way"`
	ChooseStatus string  `json:"choose_status"`
	AllWeek      string  `json:"all_week"`
	Day          int     `json:"day"`
	Session      string  `json:"session"`
	Campus       string  `json:"campus"`
	Building     string  `json:"building"`
	Classroom    string  `json:"classroom"`
	Term         int     // 学期
	Year         int
	UserID       uint
}

// AfterCreate 新建之后回调
func (s *Schedule) AfterCreate() error {
	CBNewCourseEvaluate(s.UserID, s.CourseID, s.LessonID, s.CourseName)
	return nil
}

// ScheduleList 课程表数组
type ScheduleList []Schedule

// GetSchedules 获取某个用户的某个学期课程表信息
// TODO: 待完善
func GetSchedules(userID uint, year, term int) []Schedule {
	var schedules []Schedule
	DB().Where(Schedule{
		UserID: userID,
		Term:   term,
		Year:   year,
	}).Find(&schedules)
	return schedules
}

// UpdateSchedules 更新课程表
func UpdateSchedules(userID uint, year, term int) error {
	c, err := GetJwc(userID)
	if err != nil {
		return err
	}
	defer jwc.Logout(c)

	schedules, err := schedule.Get(c)
	if err != nil {
		return err
	}

	// 删除所有的数据
	DB().Unscoped().Delete(Schedule{}, Schedule{
		UserID: userID,
		Term:   term,
		Year:   year,
	})

	for _, v := range schedules {
		newSchedule := Schedule{
			Year:         year,
			Term:         term,
			UserID:       userID,
			Project:      v.Project,
			CourseID:     v.CourseID,
			CourseName:   v.CourseName,
			LessonID:     v.LessonID,
			Credit:       v.Credit,
			CourseType:   v.CourseType,
			ExamType:     v.ExamType,
			Teachers:     strings.Join(v.Teachers, ","),
			StudyWay:     v.StudyWay,
			ChooseStatus: v.ChooseType,
			AllWeek:      v.AllWeek,
			Day:          v.Day,
			Session:      v.Session,
			Campus:       v.Campus,
			Building:     v.Building,
			Classroom:    v.Classroom,
		}

		if err := DB().Create(&newSchedule).Error; err != nil {
			return err
		}
	}
	return nil
}

// deleteMore 删除多余的数据
// 传入最新抓取的数据
func (sl ScheduleList) deleteMore(s []schedule.Schedule) error {
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

	return DB().Unscoped().Delete(Schedule{}, "id in (?)", deleteID).Error
}
