package model

// Detail 文章
type Detail struct {
	Model
	Author   string
	Title    string `gorm:"index"`
	Content  string `gorm:"type:text;"`
	URL      string
	Category string
	Tags     []Tag `gorm:"many2many:detail_tags"`
}
