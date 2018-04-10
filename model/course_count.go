package model

// CourseCount 课程统计相关的信息
// 统计方式: 评教平均分获取最近两年的平均分信息
// 点名,考核,作业 取最近一年来的众数
type CourseCount struct {
	Model
	CourseID uint    `json:"course_id"`
	FailRate float64 `json:"fail_rate"` // 挂科率
	AvgGrade float64 `json:"avg_grade"` // 平均分
	Star     float64 `json:"star"`      // 评教平均分
	Call     int     `json:"call"`      // 点名/签到方式: 0: 不点名, 1: 偶尔抽点 2: 偶尔全点 3: 全点
	ExamType int     `json:"exam_type"` // 考核方式: 0: 论文, 1: 考试, 2:大作业, 3: 其他
	Task     int     `json:"task"`      // 作业: 0: 无作业, 2: 有作业
}
