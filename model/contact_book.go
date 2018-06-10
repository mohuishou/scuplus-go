package model

type ContactBook struct {
	Model
	Title             string `json:"title"`
	Contact           string `json:"contact"`
	ContactType       string `json:"contact_type"`
	Comment           string `json:"comment"` // 备注
	ContactCategoryID uint   `json:"contact_category_id"`
}
