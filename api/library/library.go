package library

import (
	"github.com/kataras/iris"
	"github.com/mohuishou/sculibrary-go"
	"github.com/mohuishou/scuplus-go/api"
	"github.com/mohuishou/scuplus-go/middleware"
	"github.com/mohuishou/scuplus-go/model"
)

// BindLibrary 绑定图书馆账号
func BindLibrary(ctx iris.Context) {
	studentID := ctx.FormValue("student_id")
	passsword := ctx.FormValue("password")
	if studentID == "" || passsword == "" {
		api.Error(ctx, 60400, "学号或密码不能为空", nil)
		return
	}

	//验证图书馆是否可以登录
	_, err := sculibrary.NewLibrary(studentID, passsword)
	if err != nil {
		api.Error(ctx, 60401, err.Error(), nil)
	}

	uid := middleware.GetUserID(ctx)
	lib := model.UserLibrary{
		StudentID: studentID,
		Password:  passsword,
		Verify:    1,
		UserID:    uid,
	}

	if err := model.DB().Create(&lib).Error; err != nil {
		api.Error(ctx, 60500, err.Error(), nil)
		return
	}

	api.Success(ctx, "绑定成功", nil)

}

// Loan 续借
func Loan(ctx iris.Context) {
	uid := middleware.GetUserID(ctx)
	userLibrary := model.UserLibrary{}
	model.DB().Where("user_id = ?", uid).Find(&userLibrary)
	lib, err := userLibrary.GetLibrary()
	if err != nil {
		api.Error(ctx, 600401, err.Error(), map[string]interface{}{
			"verify": userLibrary.Verify,
		})
		return
	}

	bookID := ctx.FormValue("book_id")
	if lib.Loan(bookID) {
		api.Success(ctx, "续借成功", nil)
	} else {
		api.Error(ctx, 60002, "续借失败", nil)
	}
}

// GetBook 更新借阅的书籍
func GetBook(ctx iris.Context) {
	uid := middleware.GetUserID(ctx)
	userLibrary := model.UserLibrary{}
	model.DB().Where("user_id = ?", uid).Find(&userLibrary)
	lib, err := userLibrary.GetLibrary()
	if err != nil {
		api.Error(ctx, 60401, err.Error(), map[string]interface{}{
			"verify": userLibrary.Verify,
		})
		return
	}

	isHistory := "1"

	books := []sculibrary.LoanBook{}
	if isHistory == "1" {
		books = lib.GetLoanAll()
	} else {
		books = lib.GetLoan()
	}
	api.Success(ctx, "获取成功", books)
}
