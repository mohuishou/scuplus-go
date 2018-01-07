package user

import "github.com/mohuishou/scuplus-go/model"

// User 用户model
type User struct {
	model.Model
	StudentID string // 学号
	Password  string // 密码
}
