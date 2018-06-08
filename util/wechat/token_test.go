package wechat

import (
	"testing"
)

func TestGetAccessToken(t *testing.T) {
	t.Log(GetAccessToken(false))
}
