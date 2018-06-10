package model

// ContactCategory 校园通讯录分类
type ContactCategory struct {
	Model
	Name string `json:"name"`
	URL  string `json:"url"`
}
