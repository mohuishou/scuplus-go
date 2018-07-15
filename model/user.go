package model

import (
	"errors"
	"log"

	"github.com/mohuishou/scu"

	"github.com/mohuishou/scu/library"

	"github.com/gocolly/colly"

	"github.com/mohuishou/scu/jwc"

	"github.com/jinzhu/gorm"
	"github.com/mohuishou/scuplus-go/middleware"
	"github.com/mohuishou/scuplus-go/util/aes"
)

// User 用户model
type User struct {
	Model
	StudentID    string // 学号
	Password     string // 密码
	JwcStudentID string // 教务处学号
	JwcPassword  string // 教务处密码
	JwcVerify    int    // 教务处验证 0: 无法登录, 1: 正常
	Verify       int    // 统一认证门户验证: 0: 无法登录, 1: 正常
	Wechat       Wechat
	UserLibrary  UserLibrary
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

func (u *User) AfterCreate(scope *gorm.Scope) error {
	// 新增用户之后，新增用户设置
	DB().Create(&UserConfig{
		UserID: u.ID,
		Notify: NotifyAll,
	})
	return nil
}

// BeforeSave callback
func (u *User) BeforeSave(scope *gorm.Scope) error {
	if u.Password != "" {
		// 加密用户教务处密码
		password, err := aes.Encrypt([]byte(u.Password))
		if err != nil {
			log.Println("用户密码加密失败！", err, *u)
			return err
		}
		scope.SetColumn("password", password)
	}
	if u.JwcPassword != "" {
		// 加密用户教务处密码
		password, err := aes.Encrypt([]byte(u.JwcPassword))
		if err != nil {
			log.Println("用户密码加密失败！", err, *u)
			return err
		}
		scope.SetColumn("jwc_password", password)
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
	if u.JwcPassword != "" {
		password, err := aes.Decrypt(u.JwcPassword)
		if err != nil {
			log.Println("用户密码解密失败！", err, *u)
			return err
		}
		scope.SetColumn("jwc_password", password)
	}
	return nil
}

// GetCollector 获取采集器
func GetCollector(userID uint) (*colly.Collector, error) {
	user := User{}
	if err := DB().Find(&user, userID).Error; err != nil {
		return nil, err
	}
	return scu.NewCollector(user.StudentID, user.Password)
}

// GetJwc 获取教务处实例
func GetJwc(userID uint) (*colly.Collector, error) {
	user := User{}
	if err := DB().Find(&user, userID).Error; err != nil {
		return nil, err
	}
	return jwc.Login(user.JwcStudentID, user.JwcPassword)
}

// GetLibrary 获取图书馆实例
func GetLibrary(userID uint) (*library.Library, error) {
	userLib := UserLibrary{}
	if err := DB().Where("user_id = ?", userID).Find(&userLib).Error; err != nil {
		return nil, err
	}
	return library.NewLibrary(userLib.StudentID, userLib.Password)
}

func AfterUpdateBindJwc(uid uint) {
	del := DB().Unscoped().Where("user_id = ?", uid).Delete
	// 清空成绩表
	del(Grade{})
	// 清空课程表
	del(Schedule{})
	// 清空考表
	del(Exam{})
	// 清空评教表
	del(CourseEvaluate{})
	del(Evaluate{})
}

func AfterUpdateBindLibrary(uid uint) {
	DB().Unscoped().Where("user_id = ?", uid).Delete(LibraryBook{})
}

func AfterUpdateBindMy(uid uint) {
	del := DB().Unscoped().Where("user_id = ?", uid).Delete
	// 清空一卡通数据
	del(Ecard{})
	// 判断是否是研究生
	userConf := UserConfig{}
	DB().Where("user_id = ?", uid).Last(&userConf)
	if userConf.UserType == GraduateStudent {
		del(GraduateGrade{})
		del(GraduateSchedule{})
	}
}
