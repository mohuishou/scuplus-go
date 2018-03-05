package cache

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
	"github.com/mohuishou/scuplus-go/config"
)

var RedisPool *redis.Pool

// Init 初始化缓存
// 目前使用redis作为底层
func Init() {
	RedisPool = &redis.Pool{
		MaxIdle:   20,
		MaxActive: 100, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", config.Get().Redis.IP, config.Get().Redis.Port))
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

func Do(cmd string, args ...interface{}) (interface{}, error) {
	r := RedisPool.Get()
	defer r.Close()
	return r.Do(cmd, args...)
}
