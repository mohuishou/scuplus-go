package scujwc

import (
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"strings"
	"time"

	"log"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"

	"io"

	"bytes"

	"github.com/PuerkitoBio/goquery"
)

const (
	//DOMAIN 教务处ip/域名
	DOMAIN = "http://202.115.47.141"
)

//Jwc 教务处相关操作
type Jwc struct {
	uid      int
	password string
	client   http.Client
	isLogin  int //登录判断 1：已登录，0：尚未登录
}

//NewJwc 新建并初始化一个教务处对象
func NewJwc(uid int, password string) (j Jwc, err error) {
	j.password = password
	j.uid = uid
	j.initHTTP()
	err = j.login()
	if err != nil {
		return j, err
	}
	return j, nil
}

//initHTTP 初始化请求客户端
func (j *Jwc) initHTTP() {
	j.client = http.Client{}
	jar, _ := cookiejar.New(nil)
	j.client.Jar = jar
}

//login 登录教务处
func (j *Jwc) login() (err error) {

	url := DOMAIN + "/loginAction.do"
	param := "zjh=" + strconv.Itoa(j.uid) + "&mm=" + j.password

	doc, err := j.post(url, param)
	if err != nil {
		return err
	}
	errinfo := doc.Find("font[color=\"#990000\"]").Text()
	if errinfo != "" {
		j.isLogin = 0
		err := errors.New(string(errinfo))
		return err
	}
	j.isLogin = 1

	return nil
}

//Logout 退出登录
func (j *Jwc) Logout() (err error) {
	url := DOMAIN + "/logout.do"
	_, err = j.post(url, "loginType=platformLogin")
	j.isLogin = 0
	return err
}

// 发出post请求，用于教务处登录之后
func (j *Jwc) jPost(url, param string) (*goquery.Document, error) {
	if j.isLogin == 0 {
		err := errors.New("尚未登录，请先登录！")
		return nil, err
	}

	//退出登录
	defer func() {
		go j.Logout()
	}()

	return j.post(url, param)
}

//get 发出get请求
func (j *Jwc) get(url, param string) (*goquery.Document, error) {
	return j.request(url, param, "GET")
}

//post 发出post请求
func (j *Jwc) post(url, param string) (*goquery.Document, error) {
	return j.request(url, param, "POST")
}

//request 发出请求
func (j *Jwc) request(url, param, method string) (*goquery.Document, error) {

	//初始化请求
	req, err := http.NewRequest(method, url, strings.NewReader(param))
	if err != nil {
		return nil, err
	}

	//设置请求头
	req = setHeader(req)

	//设置超时时间为10s
	j.client.Timeout = time.Duration(10) * time.Second

	//发送请求
	resp, err := j.client.Do(req)
	if err != nil {
		return nil, err
	}

	//检测教务处访问是否正确
	if resp.StatusCode != 200 {
		return nil, errors.New("教务处访问错误：" + resp.Status)
	}

	//退出前关闭body
	defer resp.Body.Close()

	//编码转换
	utfBody, err := GbkToUtf8(resp.Body)
	if err != nil {
		return nil, err
	}

	// use utfBody using goquery
	doc, err := goquery.NewDocumentFromReader(utfBody)
	if err != nil {
		return nil, err
	}
	return doc, err
}

//setHeader 设置header
func setHeader(req *http.Request) *http.Request {
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.102 Safari/537.36")
	req.Header.Set("Accept", "text/javascript, text/html, application/xml, text/xml, */*")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Forwarded-For", randIP())
	return req
}

//randIP 生成随机ip地址
func randIP() (ip string) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 4; i++ {
		ip += strconv.Itoa(rand.Intn(235))
		if i != 3 {
			ip += "."
		}
	}
	return ip
}

//GbkToUtf8 编码转换
func GbkToUtf8(body io.Reader) (io.Reader, error) {
	reader := transform.NewReader(body, simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	utfBody := bytes.NewReader(d)
	return utfBody, nil
}

//Utf8ToGbk 编码转换
func Utf8ToGbk(s string) string {
	data, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(s)), simplifiedchinese.GBK.NewEncoder()))
	if err != nil {
		log.Println("编码转换错误:", err)
		return ""
	}
	return string(data)
}
