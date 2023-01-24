package picture

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/werneror/earlang/common"
	"github.com/werneror/earlang/config"
	"github.com/werneror/earlang/word"
)

type PicPicker interface {
	ID() string
	Name() string
	WordPictures(word string, number int) ([]string, error)
}

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
	err := os.MkdirAll(filepath.Join(config.PictureDir, config.PicPicker), os.ModePerm)
	if err != nil {
		log.Fatalf("failed to mkdir %v", err)
	}
}

func WordPictures(w word.Word, number int) ([]string, error) {
	query := w.GetQuery()
	picDirPath := filepath.Join(config.PictureDir, config.PicPicker, w.Key())

	dir, err := os.ReadDir(picDirPath)
	if err == nil {
		paths := make([]string, 0, number)
		for _, f := range dir {
			if f.IsDir() {
				continue
			}
			paths = append(paths, filepath.Join(picDirPath, f.Name()))
		}
		if len(paths) >= number {
			return paths[:number], nil
		}
	}

	urls, err := picker.WordPictures(query, number+5)
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(picDirPath, os.ModePerm)
	if err != nil {
		return nil, err
	}

	paths := make([]string, 0, number)
	for i, url := range urls {
		localPath := filepath.Join(picDirPath, fmt.Sprintf("%d.jpg", i))
		err = common.Download(url, localPath)
		if err != nil {
			return nil, err
		}
		paths = append(paths, localPath)
		if len(paths) >= number {
			break
		}
	}
	return paths, nil
}
