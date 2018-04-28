package model

// CourseCount 课程统计相关的信息
// 统计方式: 评教平均分获取最近两年的平均分信息
// 点名,考核,作业 取最近一年来的众数
type CourseCount struct {
	Model
	CourseID string  `json:"course_id"` // 课程号
	LessonID string  `json:"lesson_id"` // 课序号
	FailRate float64 `json:"fail_rate"` // 挂科率
	AvgGrade float64 `json:"avg_grade"` // 平均分
	Name     string  `json:"name"`      // 课程名
	Teacher  string  `json:"teacher"`   // 教师名，如果多个老师就 某某等
	Star     float64 `json:"star"`      // 评教平均分
	Campus   string  `json:"campus"`    // 校区
	Day      int     `json:"day"`       // 周几上课
	Credit   float64 `json:"credit"`    // 学分
	CallName int     `json:"call_name"` // 点名/签到方式: 1: 不点名, 2: 偶尔抽点 3: 偶尔全点 4: 全点
	ExamType int     `json:"exam_type"` // 考核方式: 1: 论文, 2: 考试, 3:大作业, 4: 其他
	Task     int     `json:"task"`      // 作业: 1: 无作业, 2: 有作业
}
