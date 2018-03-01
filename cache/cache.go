package cache

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
	"github.com/mohuishou/scuplus-go/config"
)

var Redis redis.Conn

// init 初始化缓存
// 目前使用redis作为底层
func init() {
	var err error
	Redis, err = redis.Dial("tcp", fmt.Sprintf("%s:%s", config.Get().Redis.IP, config.Get().Redis.Port))
	if err != nil {
		panic(err)
	}
}

func Do(cmd string, args ...interface{}) (interface{}, error) {
	return Redis.Do(cmd, args...)
}
