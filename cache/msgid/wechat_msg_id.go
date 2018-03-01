package msgid

import (
	"log"

	"fmt"
	"time"

	"github.com/mohuishou/scuplus-go/cache"
)

// 设置过期时间为6天
const expireTime = 3600 * 24 * 6

// Get get
func Get(uid uint) string {
	// 获取所有的keys
	v, err := cache.Do("KEYS", fmt.Sprintf("wechat.msg.id.%d.*", uid))
	if err != nil {
		log.Println("get cache code err:", err)
		return ""
	}
	keys, ok := v.([]interface{})
	if !ok {
		return ""
	}

	for _, k := range keys {
		key := string(k.([]byte))
		v, err := cache.Do("LPOP", key)
		if err != nil {
			log.Println("get cache code err:", err)
			return ""
		}
		return string(v.([]byte))
	}
	return ""
}

// Set set
func Set(uid uint, v string) error {
	_, err := cache.Do("RPUSH", getKey(uid), v)
	if err != nil {
		log.Println("set cache code err:", err)
	}
	// 设置过期时间
	_, err = cache.Do("Expire", getKey(uid), expireTime)
	if err != nil {
		log.Println("set cache code err:", err)
	}
	return err
}

func getKey(uid uint) string {
	timeStr := time.Now().Format("20060102")
	return fmt.Sprintf("wechat.msg.id.%d.%s", uid, timeStr)
}
