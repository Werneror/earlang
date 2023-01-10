package unfamiliar

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/werneror/earlang/config"
	"github.com/werneror/earlang/word"
)

var unfamiliarDataFilePath = filepath.Join(config.BaseDir, config.UnfamiliarDataFile)

type Unfamiliar struct {
	Words []word.Word `json:"words"`
}

func (u *Unfamiliar) LoadFromFile() error {
	if _, err := os.Stat(unfamiliarDataFilePath); os.IsNotExist(err) {
		u.Words = []word.Word{}
		return nil
	}
	content, err := os.ReadFile(unfamiliarDataFilePath)
	if err != nil {
		return errors.Wrapf(err, "failed to read unfamiliar data from file %s", unfamiliarDataFilePath)
	}
	err = json.Unmarshal(content, u)
	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal unfamiliar data")
	}
	return nil
}

func (u *Unfamiliar) SaveDataToFile() error {
	marshal, err := json.Marshal(u)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal unfamiliar data")
	}
	err = os.WriteFile(unfamiliarDataFilePath, marshal, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "failed to save unfamiliar data to file %s", unfamiliarDataFilePath)
	}
	return nil
}

func (u *Unfamiliar) AllWords() []word.Word {
	return u.Words
}

func (u *Unfamiliar) Add(newWord word.Word) {
	for _, w := range u.Words {
		if w.Key() == newWord.Key() {
			return
		}
	}
	u.Words = append(u.Words, newWord)
	err := u.SaveDataToFile()
	if err != nil {
		logrus.Errorf("failed to add word %v: %v", newWord, err)
	}
}

func NewUnfamiliar() (*Unfamiliar, error) {
	u := &Unfamiliar{}
	err := u.LoadFromFile()
	if err != nil {
		return nil, err
	}
	return u, nil
}
