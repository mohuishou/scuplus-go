package cos

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/url"
	"time"

	"github.com/mohuishou/scuplus-go/config"
)

const expiredInterval = 1000

var conf = config.Get().COS

func orignalSign() string {
	now := time.Now().Unix()
	rand.Seed(int64(now))
	rnd := rand.Int31()
	f := url.PathEscape("/" + conf.AppID + "/" + conf.Bucket + "/")
	sign := fmt.Sprintf("a=%s&k=%s&e=%d&t=%d&r=%d&f=%s&b=%s",
		conf.AppID,
		conf.SecretID,
		now+expiredInterval,
		now,
		rnd,
		f,
		conf.Bucket,
	)
	return sign
}

func Sign() string {
	origSign := orignalSign()
	h := hmac.New(sha1.New, []byte(config.Get().Youtu.SecretKey))
	h.Write([]byte(origSign))
	hm := h.Sum(nil)
	//attach orig_sign to hm
	dstSign := []byte(string(hm) + origSign)
	b64 := base64.StdEncoding.EncodeToString(dstSign)
	return b64
}
