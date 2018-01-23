package scujwc

import (
	"reflect"

	"strings"
)

//UserInfo 用户信息
type UserInfo struct {
	UID       string `json:"uid"`
	Name      string `json:"name"`
	NameEN    string `json:"name_en"`
	ID        string `json:"id"`
	Sex       string `json:"sex"`
	Type      string `json:"type"`
	Status    string `json:"status"`
	Nation    string `json:"nation"`
	Native    string `json:"native"`
	Birth     string `json:"birth"`
	Political string `json:"political"`
	College   string `json:"college"`
	Major     string `json:"major"`
	Year      string `json:"year"`
	Class     string `json:"class"`
	Campus    string `json:"campus"`
}

//UserInfo 获取用户信息
func (j *Jwc) UserInfo() (info UserInfo, err error) {
	return UserInfo{}.get(j)
}

func (u UserInfo) get(j *Jwc) (info UserInfo, err error) {
	url := DOMAIN + "/xjInfoAction.do"
	doc, err := j.jPost(url, "oper=xjxx")
	if err != nil {
		return info, err
	}

	userinfo := &UserInfo{}
	v := reflect.ValueOf(userinfo)
	elem := v.Elem()

	//对应关系
	eq := []int{0, 1, 3, 5, 6, 7, 9, 11, 12, 13, 14, 25, 26, 28, 29, 32}

	for i := 0; i < elem.NumField(); i++ {
		s := doc.Find("table td[width=\"275\"]").Eq(eq[i])
		elem.Field(i).SetString(strings.TrimSpace(s.Text()))
	}
	return *userinfo, nil
}
