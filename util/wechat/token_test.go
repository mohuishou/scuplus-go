package wechat

import (
	"testing"

	"github.com/mohuishou/scuplus-go/cache"
)

func TestGetAccessToken(t *testing.T) {
	cache.Init()
	t.Log(GetAccessToken(false))
}
