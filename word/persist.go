package word

import (
	"earlang/word/group"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func LoadWordsFromFile(filePath string) ([]group.Word, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	s := strings.Split(string(content), "\n")
	words := make([]group.Word, 0, len(s))
	for _, w := range s {
		w = strings.TrimSpace(w)
		if w == "" {
			continue
		}
		pieces := strings.SplitN(w, ",", 2)
		newWord := group.Word{
			English: strings.TrimSpace(pieces[0]),
		}
		if len(pieces) > 1 {
			newWord.Chinese = strings.TrimSpace(pieces[1])
		}
		words = append(words, newWord)
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
