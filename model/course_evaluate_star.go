package model

// CourseEvaluateStar 评价点赞表
// 对于同一用户只能对一条记录加分或者减分一次，不能多次点赞
type CourseEvaluateStar struct {
	Model
	UserID           uint `json:"user_id"`
	CourseEvaluateID uint `json:"course_evaluate_id"`
	Score            int  `json:"score"` // +1 或者 -1
}
