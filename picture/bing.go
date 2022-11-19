package picture

import (
	"earlang/common"
	"fmt"
	"net/url"

	"github.com/PuerkitoBio/goquery"
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
	picUrls := make([]string, 0, number)
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("src")
		if exists {
			picUrls = append(picUrls, href)
		}
	})
	return picUrls, nil
}
