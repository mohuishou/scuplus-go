package model

// CourseEvaluate 评价表，包含用户的评价
type CourseEvaluate struct {
	Model
	NickName   string `json:"nick_name"`
	Avatar     string `json:"avatar"`
	UserID     uint   `json:"user_id"`
	CourseID   string `json:"course_id"`               // 课程号
	LessonID   string `json:"lesson_id"`               // 课序号
	Comment    string `json:"comment"`                 // 评价信息
	CallName   int    `json:"call_name"`               // 点名/签到方式: 0: 不点名, 1: 偶尔抽点 2: 偶尔全点 3: 全点
	ExamType   int    `json:"exam_type"`               // 考核方式: 0: 论文, 1: 考试, 2:大作业, 3: 其他
	Task       int    `json:"task"`                    // 作业: 0: 无作业, 2: 有作业
	Star       int    `json:"star"`                    // 评分 1-3分
	Score      int    `json:"score"`                   // 计分, 分数代表本条评价的权重
	Status     int    `json:"status" gorm:"default:1"` // 是否显示，-1: 不显示, 1: 显示
	CourseName string `json:"course_name"`             // 课程名
}

// CBNewCourseEvaluate 回调函数：如果评教信息不存在则新建一条
// 应用于成绩新增，课程表更新，评教更新
func CBNewCourseEvaluate(uid uint, courseID, lessonID, courseName string) {
	e := CourseEvaluate{}
	DB().FirstOrInit(&e, CourseEvaluate{
		UserID:   uid,
		CourseID: courseID,
		LessonID: lessonID,
	})
	if e.ID == 0 {
		e.Status = -1
		e.CourseName = courseName
		DB().Create(&e)
	}
}
