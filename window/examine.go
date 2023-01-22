package window

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/werneror/earlang/common"
	"github.com/werneror/earlang/config"
	"github.com/werneror/earlang/examine"
	"github.com/werneror/earlang/pronunciation"
	"github.com/werneror/earlang/resource"
	"github.com/werneror/earlang/unfamiliar"
	"github.com/werneror/earlang/word"
)

type examineWindow struct {
	window               fyne.Window
	currentWord          word.Word
	currentWordIndex     int
	currentCorrectOption int
	imagesGrid           *fyne.Container
	examineData          *examine.Data
	processBar           *widget.ProgressBar
	closed               bool
	soundLock            sync.Mutex
}

func (e *examineWindow) showError(err error) {
	dialog.ShowError(err, e.window)
}

func (e *examineWindow) showImages(imagePaths []string) {
	e.imagesGrid.RemoveAll()
	for _, pic := range imagePaths {
		image := canvas.NewImageFromFile(pic)
		image.FillMode = canvas.ImageFillContain
		image.SetMinSize(fyne.Size{
			Width:  config.PicMinWidth,
			Height: config.PicMinHeight,
		})
		e.imagesGrid.Add(image)
	}
	e.imagesGrid.Refresh()
}

func (e *examineWindow) updateProcessBar() {
	c, a := e.examineData.Process()
	e.processBar.Min = 0
	e.processBar.Max = float64(a)
	e.processBar.SetValue(float64(c))
}

func (e *examineWindow) selectOption(i int) {
	logrus.Debugf("examine select option %d", i)
	wrong := false
	if i == e.currentCorrectOption {
		e.examineData.Correct(e.currentWordIndex)
	} else {
		wrong = true
		e.examineData.Wrong(e.currentWordIndex)
	}
	e.updateProcessBar()
	e.nextWord()
	if wrong {
		e.playWrongTone()
	}
}

func (e *examineWindow) playWrongTone() {
	e.soundLock.Lock()
	defer e.soundLock.Unlock()
	err := common.PlayMP3(config.WrongTonePath)
	if err != nil {
		e.showError(errors.Wrapf(err, "failed to play wrong tone"))
	}
}

func (e *examineWindow) readWord(word string) {
	if word == "" {
		return
	}
	e.soundLock.Lock()
	defer e.soundLock.Unlock()
	err := pronunciation.ReadOneWord(word)
	if err != nil {
		e.showError(errors.Wrapf(err, "failed to read word %s", word))
	}
}

func (e *examineWindow) autoReadWord() {
	for {
		if e.closed {
			break
		}
		e.readWord(e.currentWord.English)
		if config.WordReadAutoInterval <= 0 {
			logrus.Warnf("word.read_auto_interval is an invalid value %d", config.WordReadAutoInterval)
			config.WordReadAutoInterval = 2
		}
		time.Sleep(time.Second * time.Duration(config.WordReadAutoInterval))
	}
}

func (e *examineWindow) nextWord() {
	e.currentWord, e.currentWordIndex = e.examineData.SelectWord()
	logrus.Debugf("examine current word is %s, index is %d", e.currentWord.Key(), e.currentWordIndex)
	if e.currentWordIndex == -1 {
		dialog.ShowInformation("Congratulations!", "You have passed the examine.", e.window)
		time.Sleep(time.Second * 3)
		e.window.Close()
	}
	wordPicPath, interferePicPaths, err := examine.SelectPicture(e.currentWord, config.WordExamineOptionsCount-1)
	if err != nil {
		e.showError(err)
	}
	e.currentCorrectOption = rand.Intn(len(interferePicPaths) + 1)
	logrus.Debugf("examine current correct options is %d", e.currentCorrectOption)
	picPaths := append(interferePicPaths[:e.currentCorrectOption], append([]string{wordPicPath}, interferePicPaths[e.currentCorrectOption:]...)...)
	e.showImages(picPaths)
}

func (e *examineWindow) Show() {
	e.window.RequestFocus()
	e.window.CenterOnScreen()
	e.window.Show()
}

func newExamineWindow(app fyne.App, u *unfamiliar.Unfamiliar, l *word.List) (*examineWindow, error) {
	e := &examineWindow{
		window:     app.NewWindow("EarLang Examine"),
		imagesGrid: container.New(layout.NewGridLayout(config.WordExamineOptionsCount)),
	}
	e.window.SetIcon(resource.EarIcoe)

	var err error
	e.examineData, err = examine.NewExamineData(u, l)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create examine data")
	}
	if len(e.examineData.Words) < 5 {
		return nil, errors.New("there are too few words to take the examine")
	}

	e.window.SetOnClosed(func() {
		e.closed = true
		err := e.examineData.SaveExamineDataToFile()
		if err != nil {
			logrus.Errorf("failed to save examine: %v", err)
		}
	})

	optionButtons := make([]fyne.CanvasObject, 0, config.WordExamineOptionsCount)
	for i := 0; i < config.WordExamineOptionsCount; i++ {
		optionButtons = append(optionButtons, func(index int) *widget.Button {
			return widget.NewButton(fmt.Sprintf("%d", index+1), func() {
				e.selectOption(index)
			})
		}(i))
	}
	optionsGrid := container.New(layout.NewGridLayout(config.WordExamineOptionsCount), optionButtons...)

	e.window.Canvas().SetOnTypedKey(func(event *fyne.KeyEvent) {
		option, err := strconv.Atoi(string(event.Name))
		if err == nil && option >= 1 && option <= config.WordExamineOptionsCount {
			e.selectOption(option - 1)
		}
	})

	e.processBar = widget.NewProgressBar()
	e.processBar.TextFormatter = func() string {
		return fmt.Sprintf("%.0f of %.0f", e.processBar.Value, e.processBar.Max)
	}

	bottom := container.New(layout.NewVBoxLayout(),
		optionsGrid,
		e.processBar,
	)

	e.window.SetContent(
		container.New(layout.NewBorderLayout(e.imagesGrid, bottom, nil, nil), e.imagesGrid, bottom),
	)

	go e.autoReadWord()
	e.updateProcessBar()
	e.nextWord()
	return e, nil
}
