package scujwc

import (
	"bytes"
	"time"

	"strconv"
	"strings"

	"github.com/jordic/goics"
)

//校区
const (
	//JA 江安校区
	JA = 1
	//WJ 望江校区
	WJ = 2
)

//Ical 课程表转日历
type Ical struct {
	c      goics.Componenter
	campus int
	weeks  []string
	day    int
}

//EmitICal 1
func (ical *Ical) EmitICal() goics.Componenter {
	return ical.c
}

//init 初始化设置ical头
func (ical *Ical) init() {
	c := ical.c
	c.SetType("VCALENDAR")
	c.AddProperty("PRODID", "-//Mohuishou//SCUPLUS//FYSCU")
	c.AddProperty("CALSCALE", "GREGORIAN")
	c.AddProperty("VERSION", "2.0")
	c.AddProperty("X-WR-CALNAME", "SCUPLUS-课表")
	c.AddProperty("X-WR-TIMEZONE", "Asia/Shanghai")
	c.AddProperty("METHOD", "PUBLISH")
}

func (ical *Ical) getWeeks(s Schedule) {
	ical.weeks = strings.Split(s.AllWeek, ",")
}

func (ical *Ical) getDay(s Schedule) {
	ical.day, _ = strconv.Atoi(s.Day)
}

//Campus 设置学校
func (ical *Ical) Campus(campus int) bool {
	if campus < 1 || campus > 2 {
		return false
	}
	ical.campus = campus
	return true
}

//event 生成事件
func (ical *Ical) event(s Schedule) {
	ical.getWeeks(s)
	ical.getDay(s)
	startWeek, _ := strconv.Atoi(ical.weeks[0])
	sessions := strings.Split(s.Session, ",")
	startSession, _ := strconv.Atoi(sessions[0])
	endSession, _ := strconv.Atoi(sessions[len(sessions)-1])

	//初始化时间
	cn, _ := time.LoadLocation("Asia/Chongqing")
	tm := time.Date(2017, time.February, 26, 0, 0, 0, 0, cn)
	startDay := weekTime(startWeek, ical.day)

	//teachers
	teachers := ""
	for i := range s.Teachers {
		teachers = teachers + s.Teachers[i] + " "
	}

	//事件设置
	e := goics.NewComponent()
	e.SetType("VEVENT")
	e.AddProperty("DESCRIPTION", "课程号:"+s.CourseID+"  课序号:"+s.LessonID+"  学分:"+s.Credit+" 教师："+teachers)
	e.AddProperty("LOCATION", "@"+s.Campus+"-"+s.Building+"-"+s.Classroom)
	e.AddProperty("SUMMARY", s.CourseName+"-"+s.CourseType)
	startTime := ical.eventTime(tm, startSession, startDay, 0)
	e.AddProperty("DTSTART", startTime)
	e.AddProperty("DTEND", ical.eventTime(tm, endSession, startDay, 1))
	e.AddProperty("RRULE", ical.evetRule())
	e.AddProperty("DTSTAMP", startTime)
	e.AddProperty("CREATED", startTime)
	e.AddProperty("LAST-MODIFIED", startTime+"Z")
	e.AddProperty("SEQUENCE", "0")
	e.AddProperty("STATUS", "CONFIRMED")
	e.AddProperty("TRANSP", "OPAQUE")

	ical.c.AddComponent(e)
}

//Bytes 以[]byte返日历信息
func (ical *Ical) Bytes() []byte {
	return ical.writer().Bytes()
}

//String 返回字符串格式
func (ical *Ical) String() string {
	return ical.writer().String()
}

func (ical *Ical) writer() *bytes.Buffer {
	w := &bytes.Buffer{}
	enc := goics.NewICalEncode(w)
	enc.Encode(ical)
	return w
}

//Calendar 生成日历
func (j *Jwc) Calendar(campus int) (ical Ical, err error) {
	schedules, err := j.Schedule()
	if err != nil {
		return ical, err
	}

	ical = Ical{c: goics.NewComponent(), campus: campus}
	ical.init()

	for i := range schedules {
		ical.event(schedules[i])
	}

	return ical, nil
}

//返回循环规则
func (ical Ical) evetRule() (rule string) {

	rule = "FREQ=WEEKLY;COUNT=" + strconv.Itoa(len(ical.weeks))
	//间隔的周次
	interval := 1
	if len(ical.weeks) > 1 {
		a1, _ := strconv.Atoi(ical.weeks[1])
		a0, _ := strconv.Atoi(ical.weeks[0])
		interval = a1 - a0
	}
	rule = rule + ";INTERVAL=" + strconv.Itoa(interval)

	if ical.day < 1 || ical.day > 7 {
		return rule
	}

	//上课的时间，周几
	byDays := [7]string{"MO", "TU", "WE", "TH", "FR", "SA", "SU"}
	rule = rule + ";BYDAY=" + byDays[ical.day-1]
	return rule
}

func (ical Ical) eventTime(t time.Time, session, day, isEnd int) string {
	t = t.AddDate(0, 0, day)
	startClass := ical.classTime(session)
	startTime := t.Add(time.Duration(startClass[0])*time.Hour + time.Duration(startClass[1])*time.Minute)
	if isEnd == 1 {
		startTime = startTime.Add(time.Duration(45) * time.Minute)
	}
	return startTime.Format("20060102T150405")
}

//weekTime 根据周次以及上课的星期几返回添加的时间，返回需要增加的天数
//当输入错误时返回-1
func weekTime(week int, day int) (days int) {
	if week < 1 || day < 1 || day > 7 {
		return -1
	}
	days = day + (week-1)*7

	if day == 7 {
		days = (week - 1) * 7
	}
	return days
}

//classTime 返回所在校区的上下课时间，返回上课时间的时，分
func (ical Ical) classTime(session int) (data [2]int) {

	if (session < 1 || session > 12) || (ical.campus < 1 || ical.campus > 2) {
		return data
	}

	//江安校区时刻表
	//只包含上课时间
	//上课时间 "0815","0910","1015","1110","1350","1445","1550","1645","1740","1920","2015","2110"
	//下课时间 "0900","0955","1100","1155","1435","1530","1635","1730","1825","2005","2100","2155"
	//下课时间=上课时间+45min
	classTimeJA := [12][2]int{
		[2]int{8, 15},
		[2]int{9, 10},
		[2]int{10, 15},
		[2]int{11, 10},
		[2]int{13, 50},
		[2]int{14, 45},
		[2]int{15, 50},
		[2]int{16, 45},
		[2]int{17, 40},
		[2]int{19, 20},
		[2]int{20, 15},
		[2]int{21, 10},
	}

	//望江校区时刻表
	//上课时间 "0800","0855","1000","1055","1400","1455","1550","1655","1750","1930","2025","2120"
	//下课时间 "0845","0940","1045","1140","1445","1540","1635","1740","1835","2015","2110","2205"
	classTimeWJ := [12][2]int{
		[2]int{8, 00},
		[2]int{8, 55},
		[2]int{10, 00},
		[2]int{10, 55},
		[2]int{14, 00},
		[2]int{14, 55},
		[2]int{15, 50},
		[2]int{16, 55},
		[2]int{17, 50},
		[2]int{19, 30},
		[2]int{20, 25},
		[2]int{21, 20},
	}

	if ical.campus == JA {
		data = classTimeJA[session-1]
	} else if ical.campus == WJ {
		data = classTimeWJ[session-1]
	}

	return data
}
