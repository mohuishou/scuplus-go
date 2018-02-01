package model

// UserInfo 用户信息
type UserInfo struct {
	Model        `json:"model"`
	EcardBalance float64 `json:"ecard_balance"` // 一卡通余额
}
