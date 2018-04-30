package model

// CourseGrade 课程成绩分数
type CourseGrade struct {
	Model
	CourseID string // 课程号
	LessonID string // 课序号
	G0       int    // 0~60人数
	G60      int    // 60~70人数
	G70      int    // 70~80人数
	G80      int    // 80~90人数
	G90      int    // 90~100人数
	Year     int    // 年份
}
