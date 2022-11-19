package window

import (
	"earlang/config"
	"earlang/picture"
	"earlang/pronunciation"
	"earlang/resource"
	"earlang/word"
	"fmt"
	"os/exec"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
)

type MainWindow struct {
	window      fyne.Window
	imagesGrid  *fyne.Container
	list        *word.List
	word        string
	wordLock    sync.RWMutex
	preloadLock sync.Mutex
}

func (m *MainWindow) showError(err error) {
	dialog.ShowError(err, m.window)
}

func (m *MainWindow) getWord() string {
	var w string
	m.wordLock.RLock()
	w = m.word
	defer m.wordLock.RUnlock()
	return w
}

func (m *MainWindow) setWord(w string) {
	m.wordLock.Lock()
	m.word = w
	m.wordLock.Unlock()
}

func (m *MainWindow) readWord() {
	err := pronunciation.ReadOneWord(m.getWord())
	if err != nil {
		m.showError(fmt.Errorf("failed to read word %s: %v", m.getWord(), err))
	}
}

func (m *MainWindow) showImages() {
	pics, err := picture.WordPictures(m.getWord(), config.PicTotalNumber)
	if err != nil {
		m.showError(fmt.Errorf("failed to download pictures for word %s: %v", m.getWord(), err))
		return
	}
	m.imagesGrid.RemoveAll()
	for _, pic := range pics {
		image := canvas.NewImageFromFile(pic)
		image.FillMode = canvas.ImageFillContain
		image.SetMinSize(fyne.Size{
			Width:  config.PicMinWidth,
			Height: config.PicMinHeight,
		})
		m.imagesGrid.Add(image)
	}
	m.imagesGrid.Refresh()
}

func (m *MainWindow) preload() {
	if !m.list.Verge() {
		return
	}
	exists, newWord := m.list.PickWord()
	if exists {
		m.preloadLock.Lock()
		defer m.preloadLock.Unlock()
		logrus.Debugf("start preload %s", newWord)
		_, _ = pronunciation.WordPron(newWord, config.PronRegion)
		_, _ = picture.WordPictures(newWord, config.PicTotalNumber)
		logrus.Debugf("finish preload %s", newWord)
	}
}

func (m *MainWindow) autoReadWord() {
	for {
		if config.WordReadMode == config.WordReadModeAuto {
			w := m.getWord()
			if w != "" {
				err := pronunciation.ReadOneWord(w)
				if err != nil {
					dialog.ShowError(fmt.Errorf("failed to read word %s: %v", m.getWord(), err), m.window)
				}
			}
		}
		if config.WordReadAutoInterval <= 0 {
			logrus.Warnf("word.read_auto_interval is an invalid value %d", config.WordReadAutoInterval)
			config.WordReadAutoInterval = 2
		}
		time.Sleep(time.Second * time.Duration(config.WordReadAutoInterval))
	}
}

func (m *MainWindow) showWord() {
	m.showImages()
	if config.WordReadMode == config.WordReadModeOnce {
		m.readWord()
	}
	go m.preload()
}

func (m *MainWindow) nextWord() bool {
	exists, newWord := m.list.NextWord()
	if exists {
		m.setWord(newWord)
	}
	return exists
}

func (m *MainWindow) prevWord() bool {
	exists, newWord := m.list.PrevWord()
	if exists {
		m.setWord(newWord)
	}
	return exists
}

func (m *MainWindow) currentWord() bool {
	exists, newWord := m.list.CurrentWord()
	if exists {
		m.setWord(newWord)
	}
	return exists
}

func (m *MainWindow) ShowAndRun() {
	m.window.ShowAndRun()
}

func NewMainWindow(app fyne.App, version string) *MainWindow {
	mainWindow := &MainWindow{
		window:      app.NewWindow("EarLang"),
		imagesGrid:  container.New(layout.NewGridLayout(config.PicNumPerLine)),
		list:        word.NewList(),
		word:        "",
		wordLock:    sync.RWMutex{},
		preloadLock: sync.Mutex{},
	}
	go mainWindow.autoReadWord()

	exists := mainWindow.currentWord()
	if !exists {
		mainWindow.showError(fmt.Errorf("failed to load the word to learn"))
	} else {
		mainWindow.showWord()
	}

	prevButton := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		logrus.Debug("click prev")
		exists := mainWindow.prevWord()
		if exists {
			mainWindow.showWord()
		}
		logrus.Debug("finish prev")
	})
	readButton := widget.NewButtonWithIcon("", theme.VolumeUpIcon(), func() {
		logrus.Debug("click read")
		mainWindow.readWord()
		logrus.Debug("finish read")
	})
	nextButton := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		logrus.Debug("click next")
		exists := mainWindow.nextWord()
		if exists {
			mainWindow.showWord()
		}
		logrus.Debug("finish next")
	})
	controlGrid := container.New(layout.NewCenterLayout(), container.New(layout.NewGridLayout(3), prevButton, readButton, nextButton))

	toolbar := widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.FolderOpenIcon(), func() {
			err := exec.Command(`explorer.exe`, config.BaseDir).Start()
			if err != nil {
				mainWindow.showError(err)
			}
		}),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			newSettingWindow(app).Show()
		}),
		widget.NewToolbarAction(theme.InfoIcon(), func() {
			newAboutWindow(app, version).Show()
		}),
	)

	mainGrid := container.New(layout.NewBorderLayout(toolbar, controlGrid, nil, nil), toolbar, controlGrid, mainWindow.imagesGrid)

	mainWindow.window.SetContent(mainGrid)
	mainWindow.window.SetIcon(resource.EarIcoe)
	mainWindow.window.SetMaster()
	return mainWindow
}
