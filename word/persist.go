package word

import (
	"earlang/word/group"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

func LoadWordsFromFile(filePath string) ([]group.Word, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	s := strings.Split(string(content), "\n")
	words := make([]group.Word, 0, len(s))
	for _, w := range s {
		if w == "" {
			continue
		}
		pieces := strings.SplitN(w, ",", 2)
		if len(pieces) != 2 {
			logrus.Warnf("invalid word %s from %s", w, filePath)
		}
		words = append(words, group.Word{English: pieces[0], Chinese: pieces[1]})
	}
	return words, nil
}

func SaveWordsToFile(filePath string, words []group.Word) error {
	s := make([]string, 0, len(words))
	for _, w := range words {
		s = append(s, fmt.Sprintf("%s,%s", w.English, w.Chinese))
	}
	err := os.WriteFile(filePath, []byte(strings.Join(s, "\n")), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func LoadPointerFromFile(filePath string) (int, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return 0, nil
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		return 0, err
	}
	p, err := strconv.Atoi(string(content))
	if err != nil {
		return 0, err
	}
	return p, nil
}

func SavePointerToFile(filePath string, p int) error {
	err := os.WriteFile(filePath, []byte(strconv.FormatInt(int64(p), 10)), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
