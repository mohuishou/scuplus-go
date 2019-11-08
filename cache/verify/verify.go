package verify

import (
	"log"
	"errors"
	"fmt"

	"github.com/mohuishou/scuplus-go/cache"
)

const expireTime = 3600 * 2

var fileds = []string{
	"my",
	"jwc",
	"library",
}

// Get get
func Get(uid uint, filed string) (bool, error) {
	v, err := cache.Do("HGET", getKey(uid), filed)
	if err != nil {
		log.Println("get cache token err:", err)
	}
	t, ok := v.([]byte)
	if !ok {
		return false, errors.New("没有获取到")
	}
	return string(t) == "1", nil
}

// Set set
func Set(uid uint, filed string, v int) error {
	if !checkFiled(filed) {
		return errors.New("不存在的filed")
	}
	key := getKey(uid)
	_, err := cache.Do("HSET", key, filed, v)
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

func checkFiled(filed string) bool {
	for _, f := range fileds {
		if f == filed {
			return true
		}
	}
	return false
}

func getKey(uid uint) string {
	return fmt.Sprintf("verify.%d", uid)
}
