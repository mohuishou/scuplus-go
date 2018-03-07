package github

import (
	"log"
	"testing"
)

func Test_test(t *testing.T) {
	log.Println(client.Issues.Get(ctx, "mohuishou", "scuplus-wechat", 1))
}
