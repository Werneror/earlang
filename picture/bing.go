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
	picUrls := make([]string, 0, number)
	for {
		baseURL := "https://cn.bing.com/images/async?q=%s&first=%d&count=%d&relp=%d&tsc=ImageBasicHover&datsrc=I&layout=RowBased&mmasync=1"
		body, err := common.ReqGET(fmt.Sprintf(baseURL, url.QueryEscape(word), pageNum*number, number, number))
		if err != nil {
			return nil, fmt.Errorf("query bing images search for %s error: %w", word, err)
		}
		defer body.Close()
		doc, err := goquery.NewDocumentFromReader(body)
		if err != nil {
			return nil, fmt.Errorf("parse bing images search for %s error: %w", word, err)
		}
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
		// 有时候一页搜索结果里凑不够指定数量的图片，需要翻页
		pageNum += 1
		// 为防止死循环，限制最大页数为 5
		if pageNum > 5 || len(picUrls) >= number {
			break
		}
	}
	return picUrls, nil
}
