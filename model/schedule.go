package model

// Schedule 课程表
type Schedule struct {
	CourseID     uint
	UserID       uint
	CourseType   string // 课程属性
	StudyWay     string // 修读方式
	ChooseStatus string // 选课状态
}
