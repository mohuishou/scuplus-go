package model

import (
	"log"
	"strconv"

	"github.com/mohuishou/scujwc-go"
	"github.com/mohuishou/scuplus-go/util/aes"
)

// User 用户model
type User struct {
	Model
	StudentID string // 学号
	Password  string // 密码
	JwcVerify int    // 教务处验证: 0: 无法登录, 1: 正常
	Wechat    Wechat
}

// BeforeSave callback
func (u *User) BeforeSave() {
	if u.Password != "" {
		// 加密用户教务处密码
		password, err := aes.Encrypt([]byte(u.Password))
		if err != nil {
			log.Println("用户密码加密失败！", err, *u)
		}
		u.Password = password
	}
}

// AfterFind callback
func (u *User) AfterFind() {
	if u.Password != "" {
		password, err := aes.Decrypt(u.Password)
		if err != nil {
			log.Println("用户密码解密失败！", err, *u)
		}
		u.Password = string(password)
	}
}

// GetJwc 获取教务处实例
func (u User) GetJwc() (*scujwc.Jwc, error) {
	sid, err := strconv.Atoi(u.StudentID)
	if err != nil {
		return nil, err
	}

	jwc, err := scujwc.NewJwc(sid, u.Password)
	if err != nil {
		return nil, err
	}
	return &jwc, nil
}

// GetJwc 获取教务处实例
func GetJwc(userID uint) (*scujwc.Jwc, error) {
	user := User{}
	if err := DB().Find(&user, userID).Error; err != nil {
		return nil, err
	}
	return user.GetJwc()
}
