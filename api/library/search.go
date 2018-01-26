package library

import (
	"log"

	"github.com/kataras/iris"
	"github.com/mohuishou/sculibrary-go"
	"github.com/mohuishou/scuplus-go/api"
)

type SearchParam struct {
	Keyword  string `form:"keyword"`
	KeyType  string `form:"key_type"`
	NextPage string `form:"next_page"`
}

// Search 图书搜索
func Search(ctx iris.Context) {
	params := SearchParam{}
	ctx.ReadForm(&params)
	log.Println(params)
	res := sculibrary.Search(params.Keyword, params.KeyType, params.NextPage)
	api.Success(ctx, "获取成功！", res)
}
