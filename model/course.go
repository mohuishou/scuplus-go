package model

// Course 课程
type Course struct {
	Model
	College         string           // 学院
	CourseID        string           // 课程号
	Name            string           // 课程名
	LessonID        string           // 课序号
	Credit          float64          // 学分
	ExamType        string           // 考试类型
	AllWeek         string           // 周次: 1,2,3,4
	Day             int              // 星期
	Session         string           // 节次 1,2
	Campus          string           // 校区
	Building        string           // 教学楼
	Classroom       string           // 教室
	Max             int              // 课容量
	StudentNo       int              // 学生数
	CourseLimit     string           // 选课限制说明
	CourseCount     CourseCount      // 课程统计信息
	CourseEvaluates []CourseEvaluate `gorm:"many2many:course_evaluates;"` // 评价
	Teachers        []Teacher        `gorm:"many2many:course_teachers;"`  // 教师
}
