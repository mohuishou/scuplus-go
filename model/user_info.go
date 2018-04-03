package model

import (
	"log"

	"github.com/mohuishou/scu/jwc"

	"github.com/mohuishou/scu/jwc/info"

	"github.com/jinzhu/gorm"
	"github.com/mohuishou/scuplus-go/util/aes"
)

// UserInfo 用户信息
type UserInfo struct {
	Model
	UserID uint `json:"user_id"`
	info.UserInfo
}

// BeforeSave callback
func (u *UserInfo) BeforeSave(scope *gorm.Scope) error {
	log.Println("user", *u)
	if u.CardID != "" {
		// 加密用户教务处密码
		cardID, err := aes.Encrypt([]byte(u.CardID))
		if err != nil {
			log.Println("用户身份加密失败！", err, *u)
			return err
		}
		scope.SetColumn("card_id", cardID)
	}
	return nil
}

// AfterFind callback
func (u *UserInfo) AfterFind(scope *gorm.Scope) error {
	if u.CardID != "" {
		cardID, err := aes.Decrypt(u.CardID)
		if err != nil {
			log.Println("用户身份解密失败！", err, *u)
			return err
		}
		scope.SetColumn("card_id", cardID)
	}
	return nil
}

func UpdateUserInfo(uid uint) error {
	c, err := GetJwc(uid)
	if err != nil {
		return err
	}
	defer jwc.Logout(c)
	data, err := info.Get(c)
	if err != nil {
		return err
	}
	userInfo := UserInfo{
		UserID: uid,
	}

	if DB().Where(&userInfo).First(&userInfo).RecordNotFound() {
		userInfo.UserInfo = data
		return DB().Create(&userInfo).Error
	}

	userInfo.UserInfo = data
	return DB().Model(&userInfo).Update(userInfo).Error
}
