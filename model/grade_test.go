package model

import (
	"testing"
)

func TestUpdateGrades(t *testing.T) {

	// a := Grade{UserID: 0}
	// b := Grade{UserID: 1}
	// c := Grade{UserID: 0}
	// as := mapset.NewSet()
	// as.Add(a)
	// bs := mapset.NewSet()
	// bs.Add(b)
	// bs.Add(c)

	// fmt.Println(bs.Difference(as))

	//grades, err := UpdateGrades(1)
	//t.Log(grades, err)
	t.Log(UpdateGraduateSchedule(1, 2017, 2))
}
