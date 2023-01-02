package examine

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/werneror/earlang/config"
)

func randomlySelectOne(dir string) (string, error) {
	d, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}
	paths := make([]string, 0)
	for _, f := range d {
		paths = append(paths, filepath.Join(dir, f.Name()))
	}
	count := len(paths)
	if count == 0 {
		return "", fmt.Errorf("directory %s is empty", dir)
	}
	if count == 1 {
		return paths[0], nil
	}
	i := rand.Intn(count)
	return paths[i], nil
}

// SelectPicture 会选取输入单词的 1 张图片，并随机选择其他 count 个单词的图片各 1 张
func SelectPicture(englishWord string, count int) (string, []string, error) {
	picDirPath := filepath.Join(config.PictureDir, config.PicPicker)
	wordPicPath, err := randomlySelectOne(filepath.Join(picDirPath, englishWord))
	if err != nil {
		return "", nil, errors.Wrapf(err, "failed to select %s piecture", englishWord)
	}
	knownWords := map[string]bool{englishWord: true}
	interferePicPaths := make([]string, 0)
	for i := 0; i < count; i++ {
		attempts := 0
	retry:
		interfereWordDir, err := randomlySelectOne(filepath.Join(picDirPath))
		if err != nil {
			return "", nil, errors.Wrapf(err, "failed to select interfere word for %s", englishWord)
		}
		interfereWord := filepath.Base(interfereWordDir)
		if _, ok := knownWords[interfereWord]; ok {
			attempts += 1
			if attempts > 5 {
				return "", nil, fmt.Errorf("failed to select interfere word for %s", englishWord)
			}
			goto retry
		}
		knownWords[interfereWord] = true
		interferePicPath, err := randomlySelectOne(interfereWordDir)
		if err != nil {
			return "", nil, errors.Wrapf(err, "failed to select interfere word %s picture for word %s piecture",
				interfereWordDir, englishWord)
		}
		interferePicPaths = append(interferePicPaths, interferePicPath)
	}
	return wordPicPath, interferePicPaths, nil
}
