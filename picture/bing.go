package picture

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/werneror/earlang/common"
)

type BingImageSearch struct {
}

func (b *BingImageSearch) ID() string {
	return "bing"
}

func (b *BingImageSearch) Name() string {
	return "Bing Image Search"
}

func (b *BingImageSearch) WordPictures(word string, number int) ([]string, error) {
	pageNum := 0
	count := number + 10 // 有时候搜索结果少于指定的搜索数量，多搜索几个
	baseURL := "https://cn.bing.com/images/async?q=%s&first=%d&count=%d&relp=%d&tsc=ImageBasicHover&datsrc=I&layout=RowBased&mmasync=1"
	body, err := common.ReqGET(fmt.Sprintf(baseURL, url.QueryEscape(word), pageNum*count, count, count))
	if err != nil {
		return nil, fmt.Errorf("query bing images search for %s error: %w", word, err)
	}
	defer body.Close()
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, fmt.Errorf("parse bing images search for %s error: %w", word, err)
	}
	picUrls := make([]string, 0, number)
	doc.Find("img").EachWithBreak(func(i int, s *goquery.Selection) bool {
		href, exists := s.Attr("src")
		if exists && strings.HasPrefix(href, "http") {
			href = strings.Split(href, "?")[0] // url 中可能用 ?w=223&h=180 指定尺寸，这里去掉
			picUrls = append(picUrls, href)
			if len(picUrls) >= number {
				return false
			}
		}
		return true
	})
	return picUrls, nil
}
