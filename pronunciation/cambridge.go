package pronunciation

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/werneror/earlang/common"
)

type CambridgeDictionary struct {
}

func (c *CambridgeDictionary) ID() string {
	return "cambridge"
}

func (c *CambridgeDictionary) Name() string {
	return "Cambridge Dictionary"
}

func (c *CambridgeDictionary) WordPron(word string, region string) (string, error) {
	regionPath := fmt.Sprintf("/%s_pron/", region) // region is us or uk
	baseURL := "https://dictionary.cambridge.org/dictionary/english-chinese-simplified/%s"
	body, err := common.ReqGET(fmt.Sprintf(baseURL, url.QueryEscape(word)))
	if err != nil {
		return "", fmt.Errorf("query cambridge dictionary for %s error: %w", word, err)
	}
	defer body.Close()
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return "", fmt.Errorf("parse cambridge dictionary for %s error: %w", word, err)
	}
	targetAudio := ""
	doc.Find("audio.hdn source").Each(func(i int, s *goquery.Selection) {
		audioUrl, exists := s.Attr("src")
		if exists && targetAudio == "" && strings.HasSuffix(audioUrl, ".mp3") && strings.Contains(audioUrl, regionPath) {
			targetAudio = "https://dictionary.cambridge.org" + audioUrl
		}
	})
	if targetAudio == "" {
		return "", fmt.Errorf("no pronunciation found for %s", word)
	}
	return targetAudio, nil
}
