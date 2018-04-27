package main

import (
	"testing"

	"github.com/mohuishou/scuplus-go/model"
)

func Test_calCourceAll(t *testing.T) {
	c := model.Course{}
	model.DB().Where("course_id = ? and lesson_id = ?", "105369010", "131").First(&c)
	calCourceAll(c)
}

func Test_calCourceGrade(t *testing.T) {
	c := model.Course{}
	model.DB().Where("course_id = ? and lesson_id = ?", "105369010", "131").First(&c)
	calCourceGrade(c)
}

func Test_calCourse(t *testing.T) {
	calCourse()
}
