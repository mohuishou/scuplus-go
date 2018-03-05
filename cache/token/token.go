// Package token 缓存微信access token
package token

import (
	"log"

	"github.com/mohuishou/scuplus-go/cache"
)

const key = "wechatAccessToken"
const expireTime = 3600 * 1.5

// Get get
func Get() string {
	v, err := cache.Do("GET", key)
	if err != nil {
		log.Println("get cache token err:", err)
	}
	t, ok := v.([]byte)
	if !ok {
		return ""
	}
	return string(t)
}

// Set set
func Set(v string) error {
	_, err := cache.Do("SET", key, v)
	if err != nil {
		log.Println("set cache token err:", err)
	}
	// 设置过期时间
	_, err = cache.Do("Expire", key, expireTime)
	if err != nil {
		log.Println("set cache code err:", err)
	}
	return err
}
