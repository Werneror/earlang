package picture

import (
	"earlang/common"
	"earlang/config"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type PicPicker interface {
	ID() string
	Name() string
	WordPictures(word string, number int) ([]string, error)
}

var picBaseDir = path.Join(config.BaseDir, "picture")
var allPickers = []PicPicker{&BingImageSearch{}}
var picker PicPicker

func init() {
	found := false
	for _, d := range allPickers {
		if d.ID() == config.PicPicker {
			found = true
			picker = d
		}
	}
	if !found {
		log.Fatalf("invalid picture picker id %s", config.PicPicker)
	}
	err := os.MkdirAll(path.Join(picBaseDir, config.PicPicker), os.ModePerm)
	if err != nil {
		log.Fatalf("failed to mkdir %v", err)
	}
}

func WordPictures(word string, number int) ([]string, error) {
	picDirPath := path.Join(picBaseDir, config.PicPicker, fmt.Sprintf("%s_%d", word, number))

	dir, err := ioutil.ReadDir(picDirPath)
	if err == nil {
		paths := make([]string, 0, number)
		for _, f := range dir {
			if f.IsDir() {
				continue
			}
			paths = append(paths, path.Join(picDirPath, f.Name()))
		}
		if len(paths) >= number {
			return paths[:number], nil
		}
	}

	urls, err := picker.WordPictures(word, number)
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(picDirPath, os.ModePerm)
	if err != nil {
		return nil, err
	}

	paths := make([]string, 0, number)
	for i, url := range urls {
		localPath := path.Join(picDirPath, fmt.Sprintf("%d.jpg", i))
		err = common.Download(url, localPath)
		if err != nil {
			return nil, err
		}
		paths = append(paths, localPath)
	}
	return paths, nil
}
