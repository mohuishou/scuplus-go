package model

import (
	"log"
	"strings"
	"time"

	"github.com/mohuishou/scu/jwc/exam"
)

// Exam 考表
type Exam struct {
	Model
	UserID     uint      `json:"user_id" gorm:"unique_index:uid_name_course"`
	Name       string    `json:"name" gorm:"unique_index:uid_name_course"`
	Campus     string    `json:"campus"`
	Building   string    `json:"building"`
	Classroom  string    `json:"classroom"`
	CourseName string    `json:"course_name" gorm:"unique_index:uid_name_course"`
	Week       string    `json:"week"`
	Day        string    `json:"day"`
	Date       string    `json:"date"`
	Time       string    `json:"time"`
	Site       string    `json:"site"`
	Card       string    `json:"card"`
	Comment    string    `json:"comment"`
	StartTime  time.Time `json:"start_time"`
}

// convertExam 转换为model
func convertExam(e exam.Exam, uid uint) Exam {
	timeStr := e.Date + " " + strings.Split(e.Time, "-")[0]
	startTime, err := time.Parse("2006-01-02 15:04", timeStr)
	if err != nil {
		log.Println("时间转换失败！")
	}
	return Exam{
		UserID:     uid,
		Name:       e.Name,
		Campus:     e.Campus,
		Building:   e.Building,
		Classroom:  e.Classroom,
		CourseName: e.CourseName,
		Week:       e.Week,
		Day:        e.Day,
		Date:       e.Date,
		Time:       e.Time,
		Site:       e.Site,
		Card:       e.Card,
		Comment:    e.Comment,
		StartTime:  startTime,
	}
}

// UpdateExam 更新考表
func UpdateExam(uid uint) error {

	// 获取最新的考表
	c, err := GetJwc(uid)
	if err != nil {
		return err
	}
	exams := exam.Get(c)

	// 获取最新一条记录
	lastExam := Exam{}
	DB().Where("user_id = ?", uid).Order("start_time desc").Last(&lastExam)

	for _, value := range exams {
		e := convertExam(value, uid)
		old := Exam{}
		DB().FirstOrCreate(&old, Exam{
			UserID:     uid,
			Name:       value.Name,
			CourseName: value.CourseName,
		})
		if err := DB().Model(&old).Updates(e).Error; err != nil {
			log.Println("[Error] 考表更新失败，数据库错误", err)
		}
	}

	return nil
}
