package word

import (
	"math/rand"
	"time"

	"github.com/werneror/earlang/config"
)

type List struct {
	group *Group
}

func (l *List) PickWord() (bool, Word) {
	var unlearned []Word
	for _, w := range l.group.GetWords() {
		found := false
		for _, l := range l.group.GetLearnedWords() {
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
		return false, Word{}
	}
	index := 0
	if config.WordSelectMode != config.WordSelectModeOrder {
		index = rand.Intn(len(unlearned))
	}
	w := unlearned[index]
	l.group.AddLearnedWord(w)
	return true, w
}

func (l *List) Verge() bool {
	return l.group.GetProcess() >= l.group.GetWordsCount()-2
}

func (l *List) CurrentWord() (bool, Word) {
	if l.group.GetLearnedWordsCount() == 0 {
		return l.PickWord()
	}
	return true, l.group.GetCurrentWord()
}

func (l *List) NextWord() (bool, Word) {
	if l.group.GetProcess() == l.group.GetWordsCount()-1 {
		return false, Word{}
	}
	if l.group.GetLearnedWordsCount() > l.group.GetProcess()+1 {
		l.group.ProcessAddOne()
		return true, l.group.GetCurrentWord()
	}
	exists, word := l.PickWord()
	if !exists {
		return false, Word{}
	} else {
		l.group.ProcessAddOne()
		return true, word
	}
}

func (l *List) PrevWord() (bool, Word) {
	err := l.group.ProcessMinusOne()
	if err != nil {
		return false, Word{}
	}
	return true, l.group.GetCurrentWord()
}

func (l *List) Reset() {
	l.group.ResetProcess()
	l.group.ResetLearnedWords()
}

func (l *List) Progress() (int, int) {
	return l.group.GetProcess() + 1, l.group.GetWordsCount()
}

func NewList(groupName string) (*List, error) {
	g, err := NewGroup(groupName)
	if err != nil {
		return nil, err
	}
	l := &List{
		group: g,
	}
	return l, nil
}

func init() {
	rand.Seed(time.Now().Unix())
}
