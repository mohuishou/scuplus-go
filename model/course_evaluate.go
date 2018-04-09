package model

// CourseEvaluate 评价表，包含用户的评价
type CourseEvaluate struct {
	Model
	UserID   uint
	CourseID string // 课程号
	LessonID string // 课序号
	Comment  string // 评价信息
	Star     int    // 评分
}
