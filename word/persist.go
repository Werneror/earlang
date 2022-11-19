package word

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func LoadWordsFromFile(filePath string) ([]string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(content), "\n"), nil
}

func SaveWordsToFile(filePath string, words []string) error {
	err := ioutil.WriteFile(filePath, []byte(strings.Join(words, "\n")), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func LoadPointerFromFile(filePath string) (int, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return 0, nil
	}
	content, err := ioutil.ReadFile(filePath)
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
	err := ioutil.WriteFile(filePath, []byte(strconv.FormatInt(int64(p), 10)), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
