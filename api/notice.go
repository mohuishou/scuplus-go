package api

import (
	"log"

	"github.com/kataras/iris"
	"github.com/mohuishou/scuplus-go/model"
)

func GetNewestNotice(ctx iris.Context) {
	notice := model.Notice{}
	err := model.DB().Select([]string{"id", "abstract"}).Where("newest = 1").Last(&notice).Error
	if err != nil {
		Error(ctx, 404, "没有最新通知", nil)
		return
	}
	Success(ctx, "获取成功", notice)
}

// GetNotices 获取公告列表，公告最多四条
func GetNotices(ctx iris.Context) {
	res := []model.Notice{}
	if err := model.DB().Where("status > -1").Limit(4).Order("status desc").Select([]string{"id", "cover"}).Find(&res).Error; err != nil {
		Error(ctx, 50001, "公告获取失败！", err.Error())
		return
	}
	Success(ctx, "获取成功", res)
}

// GetNotice 获取公告详情
func GetNotice(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil || id == 0 {
		log.Println("[Error]: id:", id, err)
		Error(ctx, 50400, "参数错误", nil)
		return
	}
	notice := model.Notice{}
	model.DB().Find(&notice, id)
	Success(ctx, "获取成功！", notice)
}
