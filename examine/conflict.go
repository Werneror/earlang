package examine

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/werneror/earlang/common"
	"github.com/werneror/earlang/config"
)

var examineConflictWordsFilePath = filepath.Join(config.ExamineDir, config.ExamineConflictWordsFile)
var initConflictWords = [][]string{
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
	{"hair", "head", "beard", "moustache"},
	{"Halloween", "pumpkin"},
	{"island", "peninsula", "beach"},
	{"drainpipe", "guttering"},
	{"fridge", "freezer"},
	{"crockery", "plate"},
	{"crockery", "bowl"},
	{"mattress", "cushion"},
	{"duvet", "bed", "sheet"},
	{"duster", "towel"},
	{"screw", "nail"},
	{"moped", "scooter"},
	{"tyre", "wheel", "spoke"},
	{"jet", "plane"},
	{"pasta", "noodles"},
	{"beer", "cider"},
	{"dill", "cumin"},
	{"running", "jogging"},
	{"ear", "earlobe"},
	{"eye", "eyebrow", "eyelash", "eyelid"},
	{"nose", "nostril"},
	{"tongue", "tooth"},
	{"ankle", "foot", "calf", "foot"},
	{"hand", "palm"},
	{"mucus", "phlegm"},
	{"mouse", "rat"},
	{"claw", "paw"},
	{"bull", "cow"},
	{"beak", "bird"},
	{"petal", "flower", "pollen"},
	{"branch", "twig"},
	{"pine", "cedar"},
}

type ConflictWords struct {
	Words [][]string
}

func (c *ConflictWords) loadFromFile() error {
	if _, err := os.Stat(examineConflictWordsFilePath); os.IsNotExist(err) {
		c.Words = [][]string{}
		return nil
	}
	content, err := os.ReadFile(examineConflictWordsFilePath)
	if err != nil {
		return errors.Wrapf(err, "failed to read examine conflict words from file %s", examineConflictWordsFilePath)
	}
	c.Words = make([][]string, 0)
	s := strings.Split(string(content), "\n")
	for _, l := range s {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		words := strings.Split(l, ",")
		c.Words = append(c.Words, words)
	}
	return nil
}

func (c *ConflictWords) saveToFile() error {
	content := ""
	for _, words := range c.Words {
		content += strings.Join(words, ",") + "\n"
	}
	err := os.WriteFile(examineConflictWordsFilePath, []byte(content), os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "failed to save examine conflict words to file %s", examineConflictWordsFilePath)
	}
	return nil
}

// merge 函数用于合并 c.Words（一般来自文件）和初始化单词列表
func (c *ConflictWords) merge() {
	uncoveredWords := make([][]string, 0)
OUTER:
	for _, initWords := range initConflictWords {
		for i := 0; i < len(c.Words); i++ {
			if common.SubSlice(initWords, c.Words[i]) {
				continue OUTER
			}
			if common.SubSlice(c.Words[i], initWords) {
				c.Words[i] = initWords
				continue OUTER
			}
		}
		uncoveredWords = append(uncoveredWords, initWords)
	}
	for _, words := range uncoveredWords {
		c.Words = append(c.Words, words)
	}
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
	err := c.loadFromFile()
	if err != nil {
		return nil, err
	}
	c.merge()
	err = c.saveToFile()
	if err != nil {
		logrus.Warnf("failed to save conflict words: %v", err)
	}
	return c, nil
}
