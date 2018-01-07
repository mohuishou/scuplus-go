package jwc

import "github.com/mohuishou/scuplus-go/model"

type CourseTeacher struct {
	model.Model
	CourseID  uint
	TeacherID uint
}
