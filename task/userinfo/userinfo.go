package main

import (
	"log"

	"github.com/mohuishou/scuplus-go/model"
)

func main() {
	// 获取所有的绑定用户
	ids := make([]uint, 0)
	err := model.DB().Table("users").Where("verify = 1").Pluck("id", &ids).Error
	if err != nil {
		panic(err)
	}
	for _, id := range ids {
		if id < 144 {
			continue
		}
		if err := model.UpdateUserInfo(id); err != nil {
			log.Printf("user_id: %d, err is %v", id, err)
			continue
		}
		log.Printf("user_id: %d, updated", id)
	}
}
