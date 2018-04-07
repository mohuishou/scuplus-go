package model

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/mohuishou/scuplus-go/util/aes"
)

// UserLibrary 绑定的图书馆账号
type UserLibrary struct {
	Model
	UserID    uint
	StudentID string
	Password  string
	Verify    int
}

// BeforeSave callback
func (u *UserLibrary) BeforeSave(scope *gorm.Scope) error {
	log.Println("lalal:", u.Password)
	if u.Password != "" {
		// 加密用户教务处密码
		password, err := aes.Encrypt([]byte(u.Password))
		if err != nil {
			log.Println("用户图书馆密码加密失败！", err, *u)
			return err
		}
		scope.SetColumn("password", password)
	}
	return nil
}

// AfterFind callback
func (u *UserLibrary) AfterFind(scope *gorm.Scope) error {
	if u.Password != "" {
		password, err := aes.Decrypt(u.Password)
		if err != nil {
			log.Println("用户图书馆密码解密失败！", err, *u)
			return err
		}
		scope.SetColumn("password", password)
	}
	return nil
}
