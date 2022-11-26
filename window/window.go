package window

import (
	"earlang/config"
	"earlang/picture"
	"earlang/pronunciation"
	"earlang/resource"
	"earlang/word"
	"earlang/word/group"
	"fmt"
	"os/exec"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
)

type MainWindow struct {
	window        fyne.Window
	imagesGrid    *fyne.Container
	list          *word.List
	english       binding.String
	englishWidget fyne.Widget
	chinese       binding.String
	chineseWidget fyne.Widget
	bottomBox     *fyne.Container
	wordLock      sync.RWMutex
	preloadLock   sync.Mutex
	readButton    *widget.Button
	processBar    *widget.ProgressBar
	autoReadPause bool
}

func (m *MainWindow) showError(err error) {
	dialog.ShowError(err, m.window)
}

func (m *MainWindow) getWord() string {
	m.wordLock.RLock()
	defer m.wordLock.RUnlock()
	w, err := m.english.Get()
	if err != nil {
		m.showError(err)
		w = ""
	}
	return w
}

func (m *MainWindow) setWord(w group.Word) {
	m.wordLock.Lock()
	m.wordLock.Unlock()
	err := m.english.Set(w.English)
	if err != nil {
		m.showError(err)
	}
	err = m.chinese.Set(w.Chinese)
	if err != nil {
		m.showError(err)
	}
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
		_, _ = pronunciation.WordPron(newWord.English, config.PronRegion)
		_, _ = picture.WordPictures(newWord.English, config.PicTotalNumber)
		logrus.Debugf("finish preload %s", newWord)
	}
}

func (m *MainWindow) autoReadWord() {
	for {
		if config.WordReadMode == config.WordReadModeAuto && !m.autoReadPause {
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
	if config.WordEnglishShow {
		m.englishWidget.Show()
	} else {
		m.englishWidget.Hide()
	}
	if config.WordChineseShow {
		m.chineseWidget.Show()
	} else {
		m.chineseWidget.Hide()
	}
	m.showImages()
	if config.WordReadMode == config.WordReadModeOnce {
		m.readWord()
	}
	p, t := m.list.Progress()
	m.processBar.Min = 0
	m.processBar.Max = float64(t)
	m.processBar.SetValue(float64(p))
	go m.preload()
}

func (m *MainWindow) tempShowEnglish() {
	if m.englishWidget.Visible() {
		return
	}
	m.englishWidget.Show()
	m.bottomBox.Refresh()
	go func() {
		time.Sleep(2 * time.Second)
		if !config.WordEnglishShow {
			m.englishWidget.Hide()
		}
	}()
}

func (m *MainWindow) tempShowChinese() {
	if m.chineseWidget.Visible() {
		return
	}
	m.chineseWidget.Show()
	m.bottomBox.Refresh()
	go func() {
		time.Sleep(2 * time.Second)
		if !config.WordChineseShow {
			m.chineseWidget.Hide()
		}
	}()
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

func (m *MainWindow) prev() {
	logrus.Debug("click prev")
	exists := m.prevWord()
	if exists {
		m.showWord()
	}
	logrus.Debug("finish prev")
}

func (m *MainWindow) read() {
	logrus.Debug("click read")
	if config.WordReadMode == config.WordReadModeAuto {
		m.autoReadPause = !m.autoReadPause
		m.updateReadButtonIcon()
	} else {
		m.readWord()
	}
	logrus.Debug("finish read")
}

func (m *MainWindow) updateReadButtonIcon() {
	if config.WordReadMode == config.WordReadModeAuto {
		if m.autoReadPause {
			m.readButton.SetIcon(theme.MediaPlayIcon())
		} else {
			m.readButton.SetIcon(theme.MediaPauseIcon())
		}
	} else {
		m.readButton.SetIcon(theme.VolumeUpIcon())
	}
	m.readButton.Refresh()
}

func (m *MainWindow) next() {
	logrus.Debug("click next")
	exists := m.nextWord()
	if exists {
		m.showWord()
	}
	logrus.Debug("finish next")
}

func (m *MainWindow) reset() {
	m.list.Reset()
	exists := m.currentWord()
	if !exists {
		m.showError(fmt.Errorf("failed to load the word to learn"))
	} else {
		m.showWord()
	}
}

func (m *MainWindow) UpdateList() error {
	list, err := word.NewList()
	if err != nil {
		return err
	}
	m.list = list
	return nil
}

func (m *MainWindow) ShowAndRun() {
	m.window.RequestFocus()
	m.window.ShowAndRun()
}

func NewMainWindow(app fyne.App, version string) *MainWindow {
	picNumPerLine := config.PicNumPerLine
	if picNumPerLine > config.PicTotalNumber {
		picNumPerLine = config.PicTotalNumber
	}

	mainWindow := &MainWindow{
		window:        app.NewWindow("EarLang"),
		imagesGrid:    container.New(layout.NewGridLayout(picNumPerLine)),
		english:       binding.NewString(),
		chinese:       binding.NewString(),
		wordLock:      sync.RWMutex{},
		preloadLock:   sync.Mutex{},
		autoReadPause: false,
	}

	prevButton := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		mainWindow.prev()
	})
	mainWindow.readButton = widget.NewButtonWithIcon("", theme.VolumeUpIcon(), func() {
		mainWindow.read()
	})
	nextButton := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		mainWindow.next()
	})

	mainWindow.englishWidget = widget.NewLabelWithData(mainWindow.english)
	mainWindow.chineseWidget = widget.NewLabelWithData(mainWindow.chinese)
	// Empty string label keeps box height unchanged when word is hidden
	wordBox := container.New(layout.NewHBoxLayout(),
		widget.NewLabel(""), mainWindow.englishWidget, mainWindow.chineseWidget, widget.NewLabel(""),
	)

	mainWindow.processBar = widget.NewProgressBar()
	mainWindow.processBar.TextFormatter = func() string {
		return fmt.Sprintf("%.0f of %.0f", mainWindow.processBar.Value, mainWindow.processBar.Max)
	}

	mainWindow.bottomBox = container.New(layout.NewVBoxLayout(),
		container.New(layout.NewCenterLayout(), wordBox),
		container.New(layout.NewCenterLayout(), container.New(layout.NewGridLayout(3), prevButton, mainWindow.readButton, nextButton)),
		mainWindow.processBar,
	)

	toolbar := widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			mainWindow.reset()
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.FolderOpenIcon(), func() {
			err := exec.Command(`explorer.exe`, config.BaseDir).Start()
			if err != nil {
				mainWindow.showError(err)
			}
		}),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			newSettingWindow(app, mainWindow).Show()
		}),
		widget.NewToolbarAction(theme.InfoIcon(), func() {
			newAboutWindow(app, version).Show()
		}),
	)

	mainGrid := container.New(
		layout.NewBorderLayout(toolbar, mainWindow.bottomBox, nil, nil),
		toolbar, mainWindow.bottomBox, mainWindow.imagesGrid,
	)

	mainWindow.window.Canvas().SetOnTypedKey(func(event *fyne.KeyEvent) {
		switch event.Name {
		case "Left", "A":
			mainWindow.prev()
		case "Up", "W":
			mainWindow.window.Clipboard().SetContent(mainWindow.getWord())
		case "Right", "D":
			mainWindow.next()
		case "Down", "S":
			mainWindow.read()
		case "E":
			mainWindow.tempShowEnglish()
		case "C":
			mainWindow.tempShowChinese()
		}
	})

	if err := mainWindow.UpdateList(); err != nil {
		mainWindow.showError(err)
		goto ShowWindow
	}

	if exists := mainWindow.currentWord(); !exists {
		mainWindow.showError(fmt.Errorf("failed to load the word to learn"))
		goto ShowWindow
	}
	mainWindow.showWord()
	mainWindow.updateReadButtonIcon()
	go mainWindow.autoReadWord()

ShowWindow:
	mainWindow.window.SetContent(mainGrid)
	mainWindow.window.SetIcon(resource.EarIcoe)
	mainWindow.window.SetMaster()
	return mainWindow
}
