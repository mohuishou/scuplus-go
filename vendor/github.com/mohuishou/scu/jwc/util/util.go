package util

import (
	"errors"
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
	// 判断是否是单双周
	if strings.Contains(w, "单") {
		return "1,3,5,7,9,11,13,15,17"
	}
	if strings.Contains(w, "双") {
		return "2,4,6,8,10,12,14,16,18"
	}

	// 去除xx周上
	w = strings.Trim(w, "周上")

	// 根据逗号分割
	strs := strings.Split(w, ",")
	for _, s := range strs {
		// 根据短横线分割
		arr := strings.Split(strings.TrimSpace(s), "-")
		if len(arr) == 1 {
			allWeek = allWeek + arr[0] + ","
		} else if len(arr) == 2 {
			start, _ := strconv.Atoi(arr[0])
			end, _ := strconv.Atoi(arr[1])
			for i := start; i <= end; i++ {
				is := strconv.Itoa(i)
				allWeek = allWeek + is + ","
			}
		}
	}
	return strings.Trim(allWeek, ",")
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
