package common

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func ReqGET(targetUrl string) (io.ReadCloser, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("query %s return %d", targetUrl, resp.StatusCode)
	}
	return resp.Body, nil
}

func Download(targetUrl string, localPath string) error {
	r, err := ReqGET(targetUrl)
	if err != nil {
		return err
	}
	defer r.Close()

	f, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	if err != nil {
		return err
	}
	return nil
}
