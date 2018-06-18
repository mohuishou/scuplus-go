package main

import "github.com/mohuishou/scuplus-go/model"

func main() {
	users := []model.User{}
	model.DB().Select([]string{"id"}).Find(&users)
	for _, v := range users {
		model.DB().FirstOrCreate(&model.UserConfig{}, model.UserConfig{
			UserID: v.ID,
			Notify: model.NotifyAll,
		})
	}
}
