package webhook

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// WebHook 接收神箭手云爬虫的webhook，并且解析数据为url.Values
func WebHook(r *http.Request) (url.Values, error) {
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	b, err = ParseGzip(b)
	if err != nil {
		return nil, err
	}
	v, err := url.ParseQuery(string(b))
	return v, err
}

// ParseGzip gzip 解压
func ParseGzip(data []byte) ([]byte, error) {
	b := new(bytes.Buffer)
	binary.Write(b, binary.LittleEndian, data)
	r, err := gzip.NewReader(b)
	if err != nil {
		log.Printf("[ParseGzip] NewReader error: %v, maybe data is ungzip", err)
		return nil, err
	} else {
		defer r.Close()
		undatas, err := ioutil.ReadAll(r)
		if err != nil {
			log.Printf("[ParseGzip]  ioutil.ReadAll error: %v", err)
			return nil, err
		}
		return undatas, nil
	}
}
