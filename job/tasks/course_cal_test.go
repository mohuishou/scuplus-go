package tasks

import (
	"testing"

	"github.com/mohuishou/scuplus-go/model"
)

func TestCalCourse(t *testing.T) {
	CalCourse(2916)
}

func Test_countGrades(t *testing.T) {
	cc := model.CourseCount{}
	model.DB().Find(&cc, 2916)
	countGrades(&cc, -1)
}
