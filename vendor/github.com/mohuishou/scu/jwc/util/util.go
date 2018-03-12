package util

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

//TeacherParse 教师解析，返回包含每个教师名字的数组
func TeacherParse(t string) (teachers []string) {
	t = strings.TrimSpace(t)
	t = strings.Replace(t, "*", "", -1)
	teachers = strings.Split(t, " ")
	return teachers
}

//WeekParse 上课时间解析
func WeekParse(w string) (allWeek string) {
	re, _ := regexp.Compile(`[1-9]\d*|单|双`)
	s := re.FindAllStringSubmatch(w, -1)
	if len(s) == 1 {
		if s[0][0] == "单" {
			return "1,3,5,7,9,11,13,15,17"
		} else if s[0][0] == "双" {
			return "2,4,6,8,10,12,14,16,18"
		}
	} else if len(s) == 2 {
		start, _ := strconv.Atoi(s[0][0])
		end, _ := strconv.Atoi(s[1][0])
		for i := start; i < end; i++ {
			is := strconv.Itoa(i)
			allWeek = allWeek + is + ","
		}
		allWeek = allWeek + s[1][0]
	} else if len(s) > 2 {
		start, _ := strconv.Atoi(s[0][0])
		end, _ := strconv.Atoi(s[1][0])
		for i := start; i < end; i++ {
			is := strconv.Itoa(i)
			allWeek = allWeek + is + ","
		}
		allWeek = allWeek + s[1][0]
		for i := 2; i < len(s); i++ {
			allWeek = allWeek + "," + s[i][0]
		}
	}
	return allWeek
}

// SessionParse 上课节次解析
func SessionParse(session string) (data string, err error) {
	session = strings.TrimSpace(session)
	sessions := strings.Split(session, "~")
	if len(sessions) != 2 {
		//todo:解析
		return "", errors.New("错误")
	}
	start, _ := strconv.Atoi(sessions[0])
	end, _ := strconv.Atoi(sessions[1])
	for i := start; i < end; i++ {
		s := strconv.Itoa(i)
		data = data + s + ","
	}
	data = data + sessions[1]
	return data, nil
}
