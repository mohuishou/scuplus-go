package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB

// Model 基本模型的定义
type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
