package model

// HelpItem 帮助
type HelpItem struct {
	Model
	Title    string `json:"title"`
	Content  string `json:"content"`
	Sort     int    `json:"sort"`
	UserName string `json:"user_name"`
}
