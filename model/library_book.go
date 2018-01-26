package model

// LibraryBook 借阅的书籍
type LibraryBook struct {
	Model
	UserID      uint
	BookID      string  `json:"book_id"` // 书籍id 用于续借
	Author      string  `json:"author"`
	Title       string  `json:"title"`
	PublishYear int     `json:"publish_year"` // 出版年
	DueDate     string  `json:"due_date"`     // 到期日期
	ReturnDate  string  `json:"return_date"`  // 归还日期(借阅历史)
	ReturnTime  string  `json:"return_time"`  // 归还时间(借阅历史)
	Arrearage   float64 `json:"arrearage"`    // 欠费
	Address     string  `json:"address"`      // 分馆
	Number      string  `json:"number"`       // 索书号(当前借阅)
}
