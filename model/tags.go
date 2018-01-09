package model

// Tag 标签
type Tag struct {
	Model
	Name string `gorm:"unique_index"`
}
