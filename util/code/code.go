package code

import (
	"math/rand"
	"time"

	cacheCode "github.com/mohuishou/scuplus-go/cache/code"
)

const strLen = 6

// New 生成验证码
// key 索引，必须具有唯一性，可以是id也可以是email或者手机号
func New(key string) (string, error) {
	// 生成随机字符串
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < strLen; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}

	// 缓存随机字符串
	err := cacheCode.Set(string(result))
	return string(result), err
}

// Check 验证码验证
func Check(key, val string) bool {
	return cacheCode.Get(key) == val
}
