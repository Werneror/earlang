package word

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/werneror/earlang/config"
)

type Word struct {
	English string
	Chinese string
}

type Group struct {
	Name                 string
	wordsFilePath        string
	Words                []Word
	learnedWordsFilePath string
	learnedWords         []Word
	learnedWordsLock     sync.Mutex
	processFilePath      string
	process              int
}

func loadWordsFromFile(filePath string) ([]Word, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	s := strings.Split(string(content), "\n")
	words := make([]Word, 0, len(s))
	for _, w := range s {
		w = strings.TrimSpace(w)
		if w == "" {
			continue
		}
		pieces := strings.SplitN(w, ",", 2)
		newWord := Word{
			English: strings.TrimSpace(pieces[0]),
		}
		if len(pieces) > 1 {
			newWord.Chinese = strings.TrimSpace(pieces[1])
		}
		words = append(words, newWord)
	}
	return words, nil
}

func saveWordsToFile(filePath string, words []Word) error {
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

func (g *Group) GetName() string {
	return g.Name
}

func (g *Group) GetWords() []Word {
	return g.Words
}

func (g *Group) GetWordsCount() int {
	return len(g.Words)
}

func (g *Group) GetCurrentWord() Word {
	return g.Words[g.process]
}

func (g *Group) LoadWordsFromFile() error {
	words, err := loadWordsFromFile(g.wordsFilePath)
	if err != nil {
		return errors.Wrapf(err, "failed to load words of group %s from file %s", g.Name, g.wordsFilePath)
	}
	g.Words = words
	return nil
}

func (g *Group) SaveWordsToFile() error {
	err := saveWordsToFile(g.wordsFilePath, g.Words)
	if err != nil {
		return errors.Wrapf(err, "failed to save words of group %s to file %s", g.Name, g.wordsFilePath)
	}
	return nil
}

func (g *Group) GetLearnedWords() []Word {
	return g.learnedWords
}

func (g *Group) GetLearnedWordsCount() int {
	return len(g.learnedWords)
}

func (g *Group) AddLearnedWord(w Word) {
	g.learnedWordsLock.Lock()
	g.learnedWords = append(g.learnedWords, w)
	g.learnedWordsLock.Unlock()
	err := g.SaveLearnedWordsToFile()
	if err != nil {
		logrus.Errorf("failed to save learned words to file: %s", err)
	}
}

func (g *Group) ResetLearnedWords() {
	g.learnedWords = []Word{}
	err := g.SaveLearnedWordsToFile()
	if err != nil {
		logrus.Errorf("failed to reset learned words: %v", err)
	}
}

func (g *Group) LoadLearnedWordsFromFile() error {
	words, err := loadWordsFromFile(g.learnedWordsFilePath)
	if err != nil {
		logrus.Debugf("failed to load learned words of group %s from file %s: %v", g.Name, g.learnedWordsFilePath, err)
		words = []Word{}
	}
	g.learnedWords = words
	return nil
}

func (g *Group) SaveLearnedWordsToFile() error {
	err := saveWordsToFile(g.learnedWordsFilePath, g.learnedWords)
	if err != nil {
		return errors.Wrapf(err, "failed to save learned words of group %s to file %s", g.Name, g.learnedWordsFilePath)
	}
	return nil
}

func (g *Group) GetProcess() int {
	return g.process
}

func (g *Group) ResetProcess() {
	g.process = 0
	err := g.SaveProcessToFile()
	if err != nil {
		logrus.Errorf("failed to reset process: %v", err)
	}
}

func (g *Group) ProcessAddOne() int {
	if g.process == len(g.Words)-1 {
		return g.process
	}
	g.process = g.process + 1
	err := g.SaveProcessToFile()
	if err != nil {
		logrus.Errorf("failed to save process to file: %v", err)
	}
	return g.GetProcess()
}

func (g *Group) ProcessMinusOne() error {
	if g.process == 0 {
		return errors.New("failed to minus process")
	}
	g.process = g.process - 1
	err := g.SaveProcessToFile()
	if err != nil {
		logrus.Errorf("failed to save process to file: %v", err)
	}
	return nil
}

func (g *Group) ProcessFileExist() bool {
	if _, err := os.Stat(g.processFilePath); os.IsNotExist(err) {
		return false
	}
	return true
}

func (g *Group) LoadProcessFromFile() error {
	if _, err := os.Stat(g.processFilePath); os.IsNotExist(err) {
		g.process = 0
		return nil
	}
	content, err := os.ReadFile(g.processFilePath)
	if err != nil {
		return errors.Wrapf(err, "failed to read process of group %s from file %s", g.Name, g.processFilePath)
	}
	p, err := strconv.Atoi(string(content))
	if err != nil {
		return errors.Wrapf(err, "failed to parse process %s of group %s from file %s", content, g.Name, g.processFilePath)
	}
	g.process = p
	return nil
}

func (g *Group) SaveProcessToFile() error {
	err := os.WriteFile(g.processFilePath, []byte(strconv.FormatInt(int64(g.process), 10)), os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "failed save process of group %s to file %s", g.Name, g.processFilePath)
	}
	return nil
}

func (g *Group) GetRealLearnedWords() []Word {
	return g.learnedWords[:g.process]
}

func NewGroup(groupName string) (*Group, error) {
	g := &Group{
		Name:                 groupName,
		wordsFilePath:        filepath.Join(config.WordDir, fmt.Sprintf("%s.list", groupName)),
		learnedWordsFilePath: filepath.Join(config.WordDir, fmt.Sprintf("%s.learned", groupName)),
		processFilePath:      filepath.Join(config.WordDir, fmt.Sprintf("%s.process", groupName)),
	}
	err := g.LoadWordsFromFile()
	if err != nil {
		return nil, err
	}
	err = g.LoadLearnedWordsFromFile()
	if err != nil {
		return nil, err
	}
	err = g.LoadProcessFromFile()
	if err != nil {
		return nil, err
	}
	return g, nil
}
