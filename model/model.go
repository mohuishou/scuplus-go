package model

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // mysql 驱动
	"github.com/mohuishou/scuplus-go/config"
)

var db *gorm.DB

// Model 基本模型的定义
type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// init 数据库初始化
func initDB() {

	// 获取配置
	conf := config.Get().Mysql

	// 初始化连接
	var err error
	db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", conf.User, conf.Password, conf.Host, conf.Port, conf.DB))
	if err != nil {
		log.Fatal("数据库连接错误：", err, conf)
	}
	db.DB().SetMaxOpenConns(5)
	//db.LogMode(true)

	autoMigrate()
}

// 自动迁移
func autoMigrate() {
	db.AutoMigrate(
		&Schedule{},
		&CourseTeacher{},
		&User{},
		&Wechat{},
		&Course{},
		&Teacher{},
		&Grade{},
		&Notice{},
		&Feedback{},
		&UserLibrary{},
		&LibraryBook{},
		&Exam{},
		&Ecard{},
		&Term{},
		&TermEvent{},
		&UserInfo{},
		&CourseGrade{},
		&CourseEvaluate{},
		&CourseEvaluateStar{},
		&CourseCount{},
		&LostFind{},
		&HelpItem{},
		&ContactCategory{},
		&ContactBook{},
		&Lecture{},
		&Evaluate{},
		&UserConfig{},
	)
}

// DB 返回db，如果不存在则初始化
func DB() *gorm.DB {
	if db == nil {
		initDB()
	}
	return db
}

// Close 关闭数据库连接
func Close() error {
	return DB().Close()
}
