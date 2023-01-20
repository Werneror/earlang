package pronunciation

import (
	"fmt"
	"net/url"
)

type YouDaoDictionary struct {
}

func (y *YouDaoDictionary) ID() string {
	return "youdao"
}

func (y *YouDaoDictionary) Name() string {
	return "YouDao Dictionary"
}

func (y *YouDaoDictionary) WordPron(word string, region string) (string, error) {
	pronType := 1
	if region == "us" {
		pronType = 2
	}
	baseURL := "https://dict.youdao.com/dictvoice?audio=%s&type=%d"
	return fmt.Sprintf(baseURL, url.QueryEscape(word), pronType), nil
}
