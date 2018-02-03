package model

// CourseTeacher 课程 - 教师 关联表
type CourseTeacher struct {
	Model
	CourseID  uint
	TeacherID uint
}
