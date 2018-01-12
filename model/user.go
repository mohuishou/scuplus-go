package model

import (
	"errors"
	"log"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/mohuishou/scujwc-go"
	"github.com/mohuishou/scuplus-go/middleware"
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

// Login 用户登录， 用户不存在就新建
func (u *User) Login() (string, error) {
	if u.Wechat.Openid == "" || u.Wechat.SessionKey == "" {
		return "", errors.New("用户信息不完整")
	}

	DB().Where("openid=?", u.Wechat.Openid).Find(&u.Wechat)
	uid := u.Wechat.UserID
	if uid == 0 {
		if err := DB().Create(u).Error; err != nil {
			return "", err
		}
		uid = u.ID
	}

	DB().Find(u, uid)
	// 生成token
	return middleware.CreateToken(uid)
}

// BeforeSave callback
func (u *User) BeforeSave(scope *gorm.Scope) error {
	log.Println("user", *u)
	if u.Password != "" {
		// 加密用户教务处密码
		password, err := aes.Encrypt([]byte(u.Password))
		if err != nil {
			log.Println("用户密码加密失败！", err, *u)
			return err
		}
		scope.SetColumn("password", password)
	}
	return nil
}

// AfterFind callback
func (u *User) AfterFind(scope *gorm.Scope) error {
	if u.Password != "" {
		password, err := aes.Decrypt(u.Password)
		if err != nil {
			log.Println("用户密码解密失败！", err, *u)
			return err
		}
		scope.SetColumn("password", password)
	}
	return nil
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
