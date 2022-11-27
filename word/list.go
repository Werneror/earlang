package word

import (
	"earlang/config"
	"earlang/word/group"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type List struct {
	groupType        string
	groupName        string
	groupFile        string
	group            group.Group
	pointer          *pointer
	learnedWordsFile string
	learnedWords     []group.Word
	learnedWordsLock sync.Mutex
}

func (l *List) PickWord() (bool, group.Word) {
	var unlearned []group.Word
	for _, w := range l.group.Words {
		found := false
		for _, l := range l.learnedWords {
			if w.English == l.English {
				found = true
				break
			}
		}
		if !found {
			unlearned = append(unlearned, w)
			if config.WordSelectMode == config.WordSelectModeOrder {
				break
			}
		}
	}
	if len(unlearned) == 0 {
		return false, group.Word{}
	}
	index := 0
	if config.WordSelectMode != config.WordSelectModeOrder {
		index = rand.Intn(len(unlearned))
	}
	w := unlearned[index]
	l.learnedWordsLock.Lock()
	l.learnedWords = append(l.learnedWords, w)
	l.learnedWordsLock.Unlock()
	err := SaveWordsToFile(l.learnedWordsFile, l.learnedWords)
	if err != nil {
		logrus.Errorf("failed to save learned word to file: %v", err)
	}
	return true, w
}

func (l *List) Verge() bool {
	return l.pointer.getValue() >= len(l.learnedWords)-2
}

func (l *List) CurrentWord() (bool, group.Word) {
	if l.pointer.getValue() < 0 {
		l.pointer.setValue(0)
	}
	if len(l.learnedWords) == 0 {
		l.pointer.setValue(0)
		return l.PickWord()
	}
	if l.pointer.getValue() >= len(l.learnedWords) {
		l.pointer.setValue(len(l.learnedWords) - 1)
	}
	return true, l.learnedWords[l.pointer.getValue()]
}

func (l *List) NextWord() (bool, group.Word) {
	l.pointer.addOne()
	if l.pointer.getValue() < len(l.learnedWords) {
		return true, l.learnedWords[l.pointer.getValue()]
	}
	exists, word := l.PickWord()
	if !exists {
		l.pointer.minusOne()
		return false, group.Word{}
	}
	return true, word
}

func (l *List) PrevWord() (bool, group.Word) {
	l.pointer.minusOne()
	if l.pointer.getValue() < len(l.learnedWords) {
		return true, l.learnedWords[l.pointer.getValue()]
	}
	return false, group.Word{}
}

func (l *List) Reset() {
	l.pointer.setValue(0)
	l.learnedWords = []group.Word{}
}

func (l *List) Progress() (int, int) {
	return l.pointer.getValue() + 1, len(l.group.Words)
}

func (l *List) loadLearnedWords() {
	if _, err := os.Stat(l.learnedWordsFile); os.IsNotExist(err) {
		l.learnedWords = []group.Word{}
	} else {
		l.learnedWords, err = LoadWordsFromFile(l.learnedWordsFile)
		if err != nil {
			logrus.Errorf("failed to load learned words from file %s: %v", l.learnedWordsFile, err)
			l.learnedWords = []group.Word{}
		}
	}
}

func (l *List) loadWordList() error {
	if l.groupType == config.WordGroupTypeBuiltin {
		for _, g := range group.Groups {
			if g.Name == l.groupName {
				l.group = g
				return nil
			}
		}
		return fmt.Errorf("invalid word group name: %s", l.groupName)
	} else {
		words, err := LoadWordsFromFile(l.groupFile)
		if err != nil {
			return fmt.Errorf("failed to load word group from file %s: %w", l.groupType, err)
		}
		l.group = group.Group{
			Name:  l.groupName,
			Words: words,
		}
		logrus.Debugf("custom words is %v", words)
	}
	return nil
}

func NewList() (*List, error) {
	groupName := config.GroupName
	if config.GroupType == config.WordGroupTypeCustom {
		filename := filepath.Base(config.GroupFile)
		ext := filepath.Ext(filename)
		groupName = strings.TrimSuffix(filename, ext)
	}

	wordDir := filepath.Join(config.BaseDir, "word")
	if _, err := os.Stat(wordDir); os.IsNotExist(err) {
		_ = os.MkdirAll(wordDir, os.ModePerm)
	}

	var groupFile string
	if filepath.IsAbs(config.GroupFile) {
		groupFile = config.GroupFile
	} else {
		groupFile = filepath.Join(config.BaseDir, config.GroupFile)
	}
	logrus.Debugf("group file is %s", groupFile)

	l := &List{
		groupType:        config.GroupType,
		groupName:        groupName,
		groupFile:        groupFile,
		learnedWordsFile: filepath.Join(wordDir, fmt.Sprintf("%s_%s", groupName, config.WordLearnedFile)),
		pointer:          newPointer(filepath.Join(wordDir, fmt.Sprintf("%s_%s", groupName, config.WordProgressFile))),
	}
	err := l.loadWordList()
	if err != nil {
		return nil, err
	}
	l.loadLearnedWords()
	return l, nil
}

func init() {
	rand.Seed(time.Now().Unix())
}
