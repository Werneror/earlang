package examine

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/werneror/earlang/config"
	"github.com/werneror/earlang/picture"
	"github.com/werneror/earlang/word"
)

var conflictWords *ConflictWords

func init() {
	var err error
	conflictWords, err = NewExamineConflictWords()
	if err != nil {
		logrus.Warnf("failed to load conflict words: %s", conflictWords)
	}
}

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
func SelectPicture(w word.Word, count int) (string, []string, error) {
	query := w.GetQuery()
	picDirPath := filepath.Join(config.PictureDir, config.PicPicker)
	picWordDirPath := filepath.Join(picDirPath, w.Key())
	if _, err := os.Stat(picWordDirPath); os.IsNotExist(err) {
		_, err := picture.WordPictures(w, 1)
		if err != nil {
			return "", nil, err
		}
	}
	wordPicPath, err := randomlySelectOne(picWordDirPath)
	if err != nil {
		return "", nil, errors.Wrapf(err, "failed to select %s(%s) piecture", w.Key(), query)
	}
	knownWords := map[string]bool{w.Key(): true}
	interferePicPaths := make([]string, 0)
	for i := 0; i < count; i++ {
		attempts := 0
	retry:
		interfereWordDir, err := randomlySelectOne(filepath.Join(picDirPath))
		if err != nil {
			return "", nil, errors.Wrapf(err, "failed to select interfere word for %s(%s)", w.Key(), query)
		}
		interfereWord := filepath.Base(interfereWordDir)
		interfereEnglishWord := strings.SplitN(interfereWord, ",", 2)[0]
		if conflictWords != nil &&
			conflictWords.Conflict(w.English, interfereEnglishWord) {
			attempts += 1
			if attempts > 5 {
				return "", nil, fmt.Errorf("no non-conflicting words found for %s(%s)", w.Key(), query)
			}
			goto retry
		}
		if _, ok := knownWords[interfereWord]; ok {
			attempts += 1
			if attempts > 5 {
				return "", nil, fmt.Errorf("failed to select interfere word for %s(%s)", w.Key(), query)
			}
			goto retry
		}
		knownWords[interfereWord] = true
		interferePicPath, err := randomlySelectOne(interfereWordDir)
		if err != nil {
			attempts += 1
			if attempts > 5 {
				return "", nil, errors.Wrapf(err, "failed to select interfere word %s picture for word %s piecture",
					interfereWordDir, w.Key())
			}
			goto retry
		}
		interferePicPaths = append(interferePicPaths, interferePicPath)
	}
	return wordPicPath, interferePicPaths, nil
}
