package code

import (
	"log"

	"github.com/mohuishou/scuplus-go/cache"
)

// Get get
func Get(key string) string {
	v, err := cache.Redis.Do("GET", "verifyCode."+key)
	if err != nil {
		log.Println("get cache code err:", err)
	}
	return string(v.([]byte))
}

// Set set
func Set(key, v string) error {
	_, err := cache.Redis.Do("SET", "verifyCode."+key, v)
	if err != nil {
		log.Println("set cache code err:", err)
	}
	return err
}
