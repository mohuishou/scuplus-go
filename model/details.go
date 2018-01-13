package model

// Detail 文章
type Detail struct {
	Model
	Author   string `json:"author"`
	Title    string `gorm:"index" json:"title"`
	Content  string `gorm:"type:text;" json:"content"`
	URL      string `json:"url"`
	Category string `json:"category"`
	Tags     []Tag  `gorm:"many2many:detail_tags" json:"tags"`
}
