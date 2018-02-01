package library

import (
	"github.com/kataras/iris"
	"github.com/mohuishou/scuplus-go/api"
	"github.com/mohuishou/scuplus-go/middleware"
	"github.com/mohuishou/scuplus-go/model"
)

// Loan 续借
func Loan(ctx iris.Context) {
	uid := middleware.GetUserID(ctx)
	lib, err := model.GetLibrary(uid)
	if err != nil {
		api.Error(ctx, 30401, err.Error(), map[string]interface{}{
			"verify": 0,
		})
		return
	}

	bookID := ctx.FormValue("book_id")
	if lib.Loan(bookID) {
		api.Success(ctx, "续借成功", nil)
	} else {
		api.Error(ctx, 30001, "续借失败", nil)
	}
}

// GetBook 更新借阅的书籍
func GetBook(ctx iris.Context) {
	uid := middleware.GetUserID(ctx)
	isHistory, err := ctx.URLParamInt("is_history")
	if err != nil {
		api.Error(ctx, 30400, "参数错误！", nil)
	}
	books, err := model.UpdateLibraryBook(uid, isHistory)
	if err != nil {
		api.Error(ctx, 30002, err.Error(), nil)
	}
	api.Success(ctx, "获取成功", books)
}
