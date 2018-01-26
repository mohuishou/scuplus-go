package library

import (
	"github.com/kataras/iris"
	"github.com/mohuishou/sculibrary-go"
	"github.com/mohuishou/scuplus-go/api"
)

// SearchParam 搜索参数
type SearchParam struct {
	Keyword  string `form:"keyword"`   // 关键词
	KeyType  string `form:"key_type"`  // 搜索类型
	NextPage string `form:"next_page"` // 下一页链接
}

// Search 图书搜索
func Search(ctx iris.Context) {
	params := SearchParam{}
	ctx.ReadForm(&params)
	res := sculibrary.Search(params.Keyword, params.KeyType, params.NextPage)
	api.Success(ctx, "获取成功！", res)
}
