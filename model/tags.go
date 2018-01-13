package model

// Tag 标签
type Tag struct {
	Model
	Name    string   `gorm:"unique_index" json:"name"`
	Details []Detail `gorm:"many2many:detail_tags" `
}
