package model

// CourseGrade 课程成绩分数
type CourseGrade struct {
	Model
	CourseID string `json:"course_id"` // 课程号
	LessonID string `json:"lesson_id"` // 课序号
	G0       int    `json:"g0"`        // 0~60人数
	G60      int    `json:"g60"`       // 60~70人数
	G70      int    `json:"g70"`       // 70~80人数
	G80      int    `json:"g80"`       // 80~90人数
	G90      int    `json:"g90"`       // 90~100人数
	Year     int    `json:"year"`      // 年份
}
