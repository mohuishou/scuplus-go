package user

import (
	"github.com/mohuishou/scuplus-go/model"
)

// Wechat 用户微信相关的信息
type Wechat struct {
	model.Model
	Openid     string
	SessionKey string
	NickName   string
	AvatarURL  string
	Gender     string
}
