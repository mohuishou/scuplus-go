package model

// Wechat 用户微信相关的信息
type Wechat struct {
	Model
	Openid     string
	SessionKey string
	NickName   string
	AvatarURL  string
	Gender     string
	UserID     uint
}
