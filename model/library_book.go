package model

import (
	"github.com/mohuishou/scu/library"
)

// LibraryBook 借阅的书籍
type LibraryBook struct {
	Model
	UserID      uint    `json:"user_id"`
	IsHistory   int     `json:"is_history"`
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

// LibraryBooks books
type LibraryBooks []LibraryBook

// convertLibraryBook 从library 转化为model
func convertLibraryBook(uid uint, isHistory int, book library.LoanBook) LibraryBook {
	return LibraryBook{
		UserID:      uid,
		IsHistory:   isHistory,
		BookID:      book.BookID,
		Author:      book.Author,
		Arrearage:   book.Arrearage,
		Address:     book.Address,
		Title:       book.Title,
		PublishYear: book.PublishYear,
		DueDate:     book.DueDate,
		ReturnDate:  book.ReturnDate,
		ReturnTime:  book.ReturnTime,
		Number:      book.Number,
	}
}

// UpdateLibraryBook 更新借阅书籍
// 历史借阅将保留在数据库备份，当前借阅不保存，直接返回
// uid: 用户id
// isHistory: 0 当前借阅, 1 历史借阅
func UpdateLibraryBook(uid uint, isHistory int) (LibraryBooks, error) {
	lib, err := GetLibrary(uid)
	if err != nil {
		return nil, err
	}

	// 获取借阅数据
	var loanBooks []library.LoanBook
	if isHistory == 0 {
		loanBooks = lib.GetLoan()
	} else {
		loanBooks = lib.GetLoanAll()
	}

	// 转换数据类型
	books := make(LibraryBooks, len(loanBooks))
	for k, v := range loanBooks {
		books[k] = convertLibraryBook(uid, isHistory, v)

		// 增量保存历史借阅记录，遇到第一个插入失败(有唯一性索引)就停止新增
		if isHistory == 1 && err == nil {
			err = DB().Create(books[k]).Error
		}
	}

	return books, nil
}
