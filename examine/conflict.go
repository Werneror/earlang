package examine

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/werneror/earlang/config"
)

var examineConflictWordsFilePath = filepath.Join(config.BaseDir, config.ExamineConflictWordsFile)

func init() {
	// TODO: 如果用户也修改了这个文件，要如何保持用户的修改，又能更新数据
	if _, err := os.Stat(examineConflictWordsFilePath); os.IsNotExist(err) {
		c := &ConflictWords{
			Words: [][]string{
				{"mountain", "hill", "lake", "valley"},
				{"canal", "ditch"},
				{"forest", "swamp", "path"},
				{"stream", "river", "waterfall"},
				{"computer", "laptop", "tablet"},
				{"cloud", "sky"},
				{"lightning", "storm"},
				{"rain", "storm"},
				{"cardigan", "sweater"},
				{"cap", "hat"},
				{"thong", "knickers", "boxer shorts"},
				{"shoes", "trainers", "slippers"},
				{"stilettos", "sandal"},
				{"jeans", "pants"},
				{"socks", "stockings", "tights"},
				{"wellingtons", "boots"},
				{"blazer", "suit"},
				{"fish", "shark"},
				{"hair", "head"},
			},
		}
		err := c.SaveToFile()
		if err != nil {
			logrus.Errorf("failed to save init conflict words: %v", err)
		}
	}
}

type ConflictWords struct {
	Words [][]string
}

func (c *ConflictWords) LoadFromFile() error {
	if _, err := os.Stat(examineConflictWordsFilePath); os.IsNotExist(err) {
		c.Words = [][]string{}
		return nil
	}
	content, err := os.ReadFile(examineConflictWordsFilePath)
	if err != nil {
		return errors.Wrapf(err, "failed to read examine conflict words from file %s", examineConflictWordsFilePath)
	}
	err = json.Unmarshal(content, c)
	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal examine conflict words")
	}
	return nil
}

func (c *ConflictWords) SaveToFile() error {
	marshal, err := json.Marshal(c)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal examine conflict words")
	}
	err = os.WriteFile(examineConflictWordsFilePath, marshal, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "failed to save examine conflict words to file %s", examineConflictWordsFilePath)
	}
	return nil
}

// Conflict 函数检查 word1 和 word2 是否冲突
func (c *ConflictWords) Conflict(word1, word2 string) bool {
	for _, words := range c.Words {
		word1In := false
		word2In := false
		for _, w := range words {
			if w == word1 {
				word1In = true
			}
			if w == word2 {
				word2In = true
			}
		}
		if word1In && word2In {
			return true
		}
	}
	return false
}

func NewExamineConflictWords() (*ConflictWords, error) {
	c := &ConflictWords{}
	err := c.LoadFromFile()
	if err != nil {
		return nil, err
	}
	return c, nil
}
