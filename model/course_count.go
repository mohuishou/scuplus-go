package model

// CourseCount 课程统计相关的信息
// 统计方式: 评教平均分获取最近两年的平均分信息
// 点名,考核,作业 取最近一年来的众数
type CourseCount struct {
	Model    `json:"model"`
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
	CallName int     `json:"call_name"` // 点名/签到方式: 0: 不点名, 1: 偶尔抽点 2: 偶尔全点 3: 全点
	ExamType int     `json:"exam_type"` // 考核方式: 0: 论文, 1: 考试, 2:大作业, 3: 其他
	Task     int     `json:"task"`      // 作业: 0: 无作业, 2: 有作业
}
