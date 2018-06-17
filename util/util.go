package util

import (
	"time"
)

const (
	// TermAutumn 秋季学期
	TermAutumn = 0
	// TermSpring 春季学期
	TermSpring = 1
)

// GetYearTerm 获取当前的年份和学期
func GetYearTerm() (year, term int) {
	year = time.Now().Year()
	m := time.Now().Month()
	if m > 2 && m < 9 {
		return year, TermSpring
	}
	return year, TermAutumn
}
