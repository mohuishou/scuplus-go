package classroom

import (
	"io/ioutil"
	"net/http"
	"strings"
)

var client = http.Client{}

// Get 获取教室列表
func Get(room string) ([]byte, error) {
	resp, err := client.Post("http://cir.scu.edu.cn/cir/XLRoomData", "application/x-www-form-urlencoded", strings.NewReader("jxlname="+room))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
