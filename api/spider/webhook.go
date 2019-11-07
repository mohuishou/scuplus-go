package spider

import (
	"log"

	"github.com/json-iterator/go"
	"github.com/mohuishou/scuplus-go/model"

	"github.com/kataras/iris/v12"
	"github.com/mohuishou/scuplus-go/util/spider/webhook"
)

// WebHook 接收神箭手云爬虫的webhook，并且保存数据
func WebHook(ctx iris.Context) {
	// 解析webhook数据
	e, err := webhook.WebHook(ctx.Request())
	if err != nil {
		log.Println(err)
		return
	}

	detail := model.Detail{}
	// 解析数据, 将数据保存到数据库
	err = jsoniter.UnmarshalFromString(e.Get("data"), &detail)
	if err != nil {
		log.SetPrefix("err")
		log.Println(err)
		return
	}

	// 开始事务
	tx := model.DB().Begin()

	// 获取标签无则创建有则获取
	for k, tag := range detail.Tags {
		if err := tx.FirstOrCreate(&tag, tag).Error; err != nil {
			tx.Rollback()
			log.Println(err)
			return
		}
		detail.Tags[k] = tag
	}

	if err := tx.Create(&detail).Error; err != nil {
		tx.Rollback()
		log.Println(err)
		return
	}

	tx.Commit()

	ctx.JSON(e.Get("data_key"))
}
