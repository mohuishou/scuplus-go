package model

// User 用户model
type User struct {
	Model
	StudentID string // 学号
	Password  string // 密码
	JwcVerify int    // 教务处验证: 0: 无法登录, 1: 正常
	Wechat    Wechat
}
