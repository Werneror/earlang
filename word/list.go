package word

import (
	"earlang/config"
	"math/rand"
	"os"
	"path"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type List struct {
	wordListFile     string
	WordLearnedFile  string
	pointer          *pointer
	wordList         []string
	learnedWords     []string
	learnedWordsLock sync.Mutex
}

func (l *List) PickWord() (bool, string) {
	var unlearned []string
	for _, w := range l.wordList {
		found := false
		for _, l := range l.learnedWords {
			if w == l {
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
		return false, ""
	}
	index := 0
	if config.WordSelectMode != config.WordSelectModeOrder {
		index = rand.Intn(len(unlearned))
	}
	w := unlearned[index]
	l.learnedWordsLock.Lock()
	l.learnedWords = append(l.learnedWords, w)
	l.learnedWordsLock.Unlock()
	err := SaveWordsToFile(l.WordLearnedFile, l.learnedWords)
	if err != nil {
		logrus.Errorf("failed to save learned word to file: %v", err)
	}
	return true, w
}

func (l *List) Verge() bool {
	return l.pointer.getValue() >= len(l.learnedWords)-2
}

func (l *List) CurrentWord() (bool, string) {
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

func (l *List) NextWord() (bool, string) {
	l.pointer.addOne()
	if l.pointer.getValue() < len(l.learnedWords) {
		return true, l.learnedWords[l.pointer.getValue()]
	}
	return l.PickWord()
}

func (l *List) PrevWord() (bool, string) {
	l.pointer.minusOne()
	if l.pointer.getValue() < len(l.learnedWords) {
		return true, l.learnedWords[l.pointer.getValue()]
	}
	return false, ""
}

func (l *List) loadLearnedWords() {
	if _, err := os.Stat(l.WordLearnedFile); os.IsNotExist(err) {
		l.learnedWords = []string{}
	} else {
		l.learnedWords, err = LoadWordsFromFile(l.WordLearnedFile)
		if err != nil {
			logrus.Errorf("failed to load learned words from file %s: %v", l.WordLearnedFile, err)
			l.learnedWords = []string{}
		}
	}
}

func (l *List) loadWordList() {
	if _, err := os.Stat(l.wordListFile); os.IsNotExist(err) {
		l.wordList = nouns
		logrus.Infof("%s not exist, use built-in dictionary", l.wordListFile)
		err = SaveWordsToFile(l.wordListFile, l.wordList)
		if err != nil {
			logrus.Errorf("failed to save word list fo file %s: %v", l.wordListFile, err)
		}
	} else {
		l.wordList, err = LoadWordsFromFile(l.wordListFile)
		if err != nil {
			logrus.Errorf("failed to load word list from file %s: %v", l.wordListFile, err)
			l.wordList = []string{}
		}
	}
}

func NewList() *List {
	l := &List{
		wordListFile:    path.Join(config.BaseDir, config.WordListFile),
		WordLearnedFile: path.Join(config.BaseDir, config.WordLearnedFile),
		pointer:         newPointer(path.Join(config.BaseDir, config.WordProgressFile)),
	}
	l.loadWordList()
	l.loadLearnedWords()
	return l
}

func init() {
	rand.Seed(time.Now().Unix())
}
