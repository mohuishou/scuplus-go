package jwc

import (
	"github.com/mohuishou/scuplus-go/model"
)

// Schedule 课程表
type Schedule struct {
	model.Model
	CourseID     uint
	UserID       uint
	CourseType   string // 课程属性
	StudyWay     string // 修读方式
	ChooseStatus string // 选课状态
}
