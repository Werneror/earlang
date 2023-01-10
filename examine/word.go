package examine

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/werneror/earlang/config"
	"github.com/werneror/earlang/unfamiliar"
	"github.com/werneror/earlang/word"
)

var examineDataFilePath = filepath.Join(config.BaseDir, config.ExamineDataFile)

type wordResult struct {
	word.Word
	CorrectTimes uint64 `json:"correct_times"` // 以往所有测试中这个单词答对的次数
	WrongTimes   uint64 `json:"wrong_times"`   // 以往所有测试中这个单词答错的次数
	correct      bool   // 当前这次测试中这个单词是否答对了
	needExamine  bool   // 当前这次测试中是否需要测试这个单词
}

type Data struct {
	Words []*wordResult `json:"words"`
}

func (d *Data) LoadExamineDataFromFile() error {
	if _, err := os.Stat(examineDataFilePath); os.IsNotExist(err) {
		d.Words = []*wordResult{}
		return nil
	}
	content, err := os.ReadFile(examineDataFilePath)
	if err != nil {
		return errors.Wrapf(err, "failed to read examine data from file %s", examineDataFilePath)
	}
	err = json.Unmarshal(content, d)
	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal examine data")
	}
	return nil
}

func (d *Data) SaveExamineDataToFile() error {
	marshal, err := json.Marshal(d)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal examine data")
	}
	err = os.WriteFile(examineDataFilePath, marshal, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "failed to save examine data to file %s", examineDataFilePath)
	}
	return nil
}

func (d *Data) SelectWord() (word.Word, int) {
	remainingWords := map[int]*wordResult{}
	for i, w := range d.Words {
		if w.needExamine && !w.correct {
			remainingWords[i] = w
		}
	}
	// TODO: 完善挑选单词的算法
	for i, w := range remainingWords {
		return w.Word, i
	}
	return word.Word{}, -1
}

func (d *Data) Process() (int, int) {
	correctWordCount := 0
	allNeedExamineWordCount := 0
	for _, w := range d.Words {
		if w.needExamine {
			allNeedExamineWordCount += 1
		}
		if w.correct {
			correctWordCount += 1
		}
	}
	return correctWordCount, allNeedExamineWordCount
}

func (d *Data) Correct(i int) {
	d.Words[i].CorrectTimes += 1
	d.Words[i].correct = true
}

func (d *Data) Wrong(i int) {
	d.Words[i].WrongTimes += 1
}

func learnedWords() (map[string]word.Word, error) {
	groups, err := word.AllGroups()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all groups")
	}
	learnedWords := make(map[string]word.Word, 0)
	for _, g := range groups {
		for _, w := range g.GetRealLearnedWords() {
			learnedWords[w.Key()] = w
		}
	}
	return learnedWords, nil
}

func allWords() (map[string]word.Word, error) {
	groups, err := word.AllGroups()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all groups")
	}
	allWords := make(map[string]word.Word, 0)
	for _, g := range groups {
		for _, w := range g.GetWords() {
			allWords[w.Key()] = w
		}
	}
	return allWords, nil
}

func unfamiliarWords(u *unfamiliar.Unfamiliar) (map[string]word.Word, error) {
	words := make(map[string]word.Word, 0)
	for _, w := range u.AllWords() {
		words[w.Key()] = w
	}
	return words, nil
}

func NewExamineData(u *unfamiliar.Unfamiliar) (*Data, error) {
	r := &Data{}
	err := r.LoadExamineDataFromFile()
	if err != nil {
		return nil, err
	}

	var words map[string]word.Word
	switch config.ExamineMode {
	case config.ExamineModeAll:
		words, err = allWords()
	case config.ExamineModeLearned:
		words, err = learnedWords()
	case config.ExamineModeUnfamiliar:
		words, err = unfamiliarWords(u)
	default:
		return nil, fmt.Errorf("unsupport examine mode: %s", config.ExamineMode)
	}

	newWords := make([]word.Word, 0)
	for key, lw := range words {
		found := false
		for _, rw := range r.Words {
			if key == rw.Key() {
				found = true
				rw.needExamine = true
				break
			}
		}
		if !found {
			newWords = append(newWords, lw)
		}
	}
	for _, nw := range newWords {
		r.Words = append(r.Words, &wordResult{
			Word:         nw,
			CorrectTimes: 0,
			WrongTimes:   0,
			needExamine:  true,
		})
	}
	return r, nil
}
