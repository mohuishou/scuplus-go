package sculibrary

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/gocolly/colly"
)

// SearchResult 搜索结果
type SearchResult struct {
	Books   []SearchBook
	NextURL string
}

// SearchBook 搜索到的书籍
type SearchBook struct {
	Author        string
	Title         string
	Cover         string
	Press         string // 出版社
	PublishYear   string // 出版年
	Number        string // 索书号(当前借阅)
	BookAddresses []SearchBookAddress
}

// SearchBookAddress 搜索到的书籍地址
type SearchBookAddress struct {
	Address   string // 馆藏地址
	Number    string // 索书号
	ALLCount  string // 馆藏数
	LoanCount string // 可借数
}

// Search 搜索
func Search(keyword, keyType, nextURL string) SearchResult {
	c := colly.NewCollector()
	searchResult := SearchResult{}
	searchBooks := make([]SearchBook, 0)

	// 获取结果列表
	c.OnHTML("#brief table.items", func(e *colly.HTMLElement) {
		searchBook := SearchBook{}
		searchBook.Cover = strings.TrimSpace(e.DOM.Find("td.cover a img").AttrOr("src", ""))
		searchBook.Title = strings.TrimSpace(e.DOM.Find("div.itemtitle a").Text())
		searchBook.Author = strings.TrimSpace(e.DOM.Find(" table:nth-child(2) > tbody > tr:nth-child(1) > td:nth-child(2)").Text())
		searchBook.Number = strings.TrimSpace(e.DOM.Find(" table:nth-child(2) > tbody > tr:nth-child(1) > td:nth-child(4)").Text())
		searchBook.Press = strings.TrimSpace(e.DOM.Find(" table:nth-child(2) > tbody > tr:nth-child(2) > td:nth-child(2)").Text())
		searchBook.PublishYear = strings.TrimSpace(e.DOM.Find(" table:nth-child(2) > tbody > tr:nth-child(2) > td:nth-child(4)").Text())
		searchBook.BookAddresses = make([]SearchBookAddress, 0)
		e.DOM.Find("td.col2 tr").Each(func(i int, s *goquery.Selection) {

			if s.Find("td.libnname a").Text() == "" {
				return
			}
			bookAddr := SearchBookAddress{}
			bookAddr.Address = strings.TrimSpace(s.Find("td.libnname").Text())
			bookAddr.Number = strings.TrimSpace(s.Find("td.boodid").Text())
			strs := strings.Split(s.Find("td.holding").Text(), "/")
			if len(strs) == 2 {
				bookAddr.ALLCount = strings.TrimSpace(strs[0])
				bookAddr.LoanCount = strings.TrimSpace(strs[1])
			}
			searchBook.BookAddresses = append(searchBook.BookAddresses, bookAddr)
		})

		searchBooks = append(searchBooks, searchBook)
	})

	// 获取下一页地址
	c.OnHTML("div#hitnum", func(e *colly.HTMLElement) {
		r, _ := regexp.Compile(`(\d+)\s+of\s+(\d+)`)
		res := r.FindAllStringSubmatch(e.Text, -1)
		if len(res[0]) == 3 {
			last, _ := strconv.Atoi(res[0][1])
			all, _ := strconv.Atoi(res[0][2])
			if last+1 < all {
				searchResult.NextURL = fmt.Sprintf("http://opac.scu.edu.cn:8080%s?func=short-jump&jump=%d&pag=now", e.Request.URL.EscapedPath(), last+1)
			}
		}
	})

	url := nextURL
	if url == "" {
		url = fmt.Sprintf("%s?func=find-b&find_code=%s&request=%s&local_base=SCU01", getURL(), keyType, keyword)
	}
	c.Visit(url)

	searchResult.Books = searchBooks
	return searchResult
}
