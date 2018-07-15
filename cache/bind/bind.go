package bind

import (
	"log"

	"fmt"

	"time"

	"strconv"

	"github.com/mohuishou/scuplus-go/cache"
)

// Get get
func Get(uid uint, t string) int {
	v, err := cache.Do("GET", getKey(uid, t))
	if err != nil {
		log.Println("get cache code err:", err)
	}
	val, ok := v.([]byte)
	if !ok {
		return 0
	}
	n, _ := strconv.Atoi(string(val))
	return n
}

// Add +1
func Add(uid uint, t string) error {
	key := getKey(uid, t)
	_, err := cache.Do("INCR", key)
	if err != nil {
		log.Println("set cache code err:", err)
	}
	// 设置过期时间
	_, err = cache.Do("Expireat", key, getExpTime())
	if err != nil {
		log.Println("set cache code err:", err)
	}
	return err
}

func getKey(uid uint, t string) string {
	return fmt.Sprintf("bind.%s.%d", t, uid)
}

// 获取过期时间，当日23.59分
func getExpTime() int64 {
	dateStr := time.Now().Format("2006-01-02")
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", dateStr+" 23:59:59", time.Local)
	return t.Unix()
}
