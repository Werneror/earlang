package pronunciation

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/werneror/earlang/common"
	"github.com/werneror/earlang/config"
)

type PronPicker interface {
	ID() string
	Name() string
	WordPron(word string, region string) (string, error)
}

var pronBaseDir = path.Join(config.BaseDir, "pronunciation")
var allPickers = []PronPicker{&CambridgeDictionary{}, &YouDaoDictionary{}}
var picker PronPicker
var defaultPhrasePicker PronPicker = &YouDaoDictionary{}

func init() {
	found := false
	for _, p := range allPickers {
		if p.ID() == config.PronPicker {
			found = true
			picker = p
		}
	}
	if !found {
		log.Fatalf("invalid pronunciation picker id %s", config.PronPicker)
	}
	err := os.MkdirAll(path.Join(pronBaseDir, config.PronPicker), os.ModePerm)
	if err != nil {
		log.Fatalf("failed to mkdir %v", err)
	}
}

func WordPron(word string, region string) (string, error) {
	isPhrase := false
	pickerID := config.PronPicker
	if strings.Contains(word, " ") ||
		strings.Contains(word, "'") ||
		strings.Contains(word, "-") {
		pickerID = defaultPhrasePicker.ID()
		isPhrase = true
	}
	audioFilePath := path.Join(pronBaseDir, pickerID, fmt.Sprintf("%s_%s.mp3", region, word))
	_, err := os.Stat(audioFilePath)
	if err == nil {
		return audioFilePath, nil
	}

	var url string
	if isPhrase {
		url, err = defaultPhrasePicker.WordPron(word, region)
		if err != nil {
			return "", err
		}
	} else {
		url, err = picker.WordPron(word, region)
		if err != nil {
			return "", err
		}
	}

	err = common.Download(url, audioFilePath)
	if err != nil {
		return "", err
	}
	return audioFilePath, nil
}

var readLock sync.Mutex

func ReadOneWord(w string) error {
	audioPath, err := WordPron(w, config.PronRegion)
	if err != nil {
		return err
	}
	readLock.Lock()
	defer readLock.Unlock()
	return common.PlayAudio(audioPath)
}
