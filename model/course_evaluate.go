package model

// CourseEvaluate 评价表，包含用户的评价
type CourseEvaluate struct {
	Model
	UserID   uint   `json:"user_id"`
	CourseID string `json:"course_id"` // 课程号
	LessonID string `json:"lesson_id"` // 课序号
	Comment  string `json:"comment"`   // 评价信息
	Call     int    `json:"call"`      // 点名/签到方式: 0: 不点名, 1: 偶尔抽点 2: 偶尔全点 3: 全点
	ExamType int    `json:"exam_type"` // 考核方式: 0: 论文, 1: 考试, 2:大作业, 3: 其他
	Task     int    `json:"task"`      // 作业: 0: 无作业, 2: 有作业
	Star     int    `json:"star"`      // 评分 0-3分,0分不计入统计
	Score    int    `json:"score"`     // 计分, 分数代表本条评价的权重
}
