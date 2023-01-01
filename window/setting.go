package window

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/werneror/earlang/config"
	"github.com/werneror/earlang/resource"
	"github.com/werneror/earlang/word"
)

type settingWindow struct {
	window     fyne.Window
	mainWindow *MainWindow
}

func (s *settingWindow) Show() {
	s.window.RequestFocus()
	s.window.Show()
}

func getGroupDisplayName(groupName string, total, process int) string {
	return fmt.Sprintf("%s (%d/%d)", groupName, process, total)
}

func displayNameToGroupName(displayName string) string {
	z := strings.Split(displayName, " (")
	z = z[:len(z)-1]
	return strings.Join(z, " (")
}

func newSettingWindow(app fyne.App, mainWindow *MainWindow) *settingWindow {
	s := &settingWindow{
		window:     app.NewWindow("EarLang Setting"),
		mainWindow: mainWindow,
	}

	logLevelSelect := widget.NewSelect([]string{"debug", "info", "warning", "error"}, func(string) {})
	logLevelSelect.SetSelected(config.LogLevel)

	pronPickerSelect := widget.NewSelect([]string{"cambridge"}, func(string) {})
	pronPickerSelect.SetSelected(config.PronPicker)

	pronRegionSelect := widget.NewSelect([]string{"us", "uk"}, func(string) {})
	pronRegionSelect.SetSelected(config.PronRegion)

	picPickerSelect := widget.NewSelect([]string{"bing"}, func(string) {})
	picPickerSelect.SetSelected(config.PicPicker)

	picNumberSelect := widget.NewSelect([]string{"5", "10", "15"}, func(string) {})
	picNumberSelect.SetSelected(fmt.Sprintf("%d", config.PicTotalNumber))

	allWords := 0
	learnedWords := 0
	groupNames := make([]string, 0)
	var currentDisplayName string
	groups, err := word.AllGroups()
	if err != nil {
		logrus.Errorf("failed to get all groups: %v", err)
	}
	for _, g := range groups {
		allWords += g.GetWordsCount()
		progress := g.GetProcess()
		learnedWords += progress
		displayName := getGroupDisplayName(g.Name, len(g.Words), progress)
		groupNames = append(groupNames, displayName)
		if config.WordGroupName == g.Name {
			currentDisplayName = displayName
		}
	}
	groupNameSelect := widget.NewSelect(groupNames, func(s string) {})
	groupNameSelect.SetSelected(currentDisplayName)
	groupNameSelectHint := fmt.Sprintf("groups: %d, total words: %d, learned words: %d", len(groups), allWords, learnedWords)

	readModeSelect := widget.NewSelect([]string{config.WordReadModeAuto, config.WordReadModeOnce, config.WordReadModeManual}, func(string) {})
	readModeSelect.SetSelected(config.WordReadMode)

	wordSelectModeSelect := widget.NewSelect([]string{config.WordSelectModeOrder, config.WordSelectModeRandom}, func(string) {})
	wordSelectModeSelect.SetSelected(config.WordSelectMode)

	readAutoIntervalSelect := widget.NewSelect([]string{"1", "2", "3", "4", "5", "10"}, func(string) {})
	readAutoIntervalSelect.SetSelected(fmt.Sprintf("%d", config.WordReadAutoInterval))

	showEnglishCheck := widget.NewCheck("show word below the pictures", func(bool) {})
	showEnglishCheck.SetChecked(config.WordEnglishShow)

	showChineseCheck := widget.NewCheck("show Chinese translation below the pictures", func(bool) {})
	showChineseCheck.SetChecked(config.WordChineseShow)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "log level", Widget: logLevelSelect},
			{Text: "pronunciation source", Widget: pronPickerSelect},
			{Text: "pronunciation region", Widget: pronRegionSelect},
			{Text: "picture source", Widget: picPickerSelect},
			{Text: "picture number", Widget: picNumberSelect},
			{Text: "word group name", Widget: groupNameSelect, HintText: groupNameSelectHint},
			{Text: "word read mode", Widget: readModeSelect},
			{Text: "word auto read interval(s)", Widget: readAutoIntervalSelect},
			{Text: "word select mode", Widget: wordSelectModeSelect},
			{Text: "show word", Widget: showEnglishCheck},
			{Text: "show Chinese", Widget: showChineseCheck},
		},
		OnSubmit: func() {
			needUpdateList := false
			needUpdateWord := false
			needUpdateReadButtonIcon := false

			config.LogLevel = logLevelSelect.Selected
			viper.Set("log.level", config.LogLevel)

			config.PronPicker = pronPickerSelect.Selected
			viper.Set("pronunciation.picker", config.PronPicker)

			config.PronRegion = pronRegionSelect.Selected
			viper.Set("pronunciation.region", config.PronRegion)

			config.PicPicker = picPickerSelect.Selected
			viper.Set("picture.picker", config.PicPicker)

			picNumber, err := strconv.Atoi(picNumberSelect.Selected)
			if err != nil {
				logrus.Error("failed to parse picture total number %s: %v", picNumberSelect.Selected, err)
				picNumber = config.PicTotalNumber
			}
			if picNumber != config.PicTotalNumber {
				needUpdateWord = true
			}
			config.PicTotalNumber = picNumber
			viper.Set("picture.total_number", config.PicTotalNumber)

			selectedGroupName := displayNameToGroupName(groupNameSelect.Selected)
			if config.WordGroupName != selectedGroupName {
				needUpdateList = true
			}
			config.WordGroupName = selectedGroupName
			viper.Set("word.group_name", config.WordGroupName)

			if config.WordReadMode != readModeSelect.Selected {
				needUpdateReadButtonIcon = true
			}
			config.WordReadMode = readModeSelect.Selected
			viper.Set("word.read_mode", config.WordReadMode)

			readAutoInterval, err := strconv.Atoi(readAutoIntervalSelect.Selected)
			if err != nil {
				logrus.Error("failed to parse word read auto interval %s: %v", readAutoIntervalSelect.Selected, err)
				readAutoInterval = config.WordReadAutoInterval
			}
			config.WordReadAutoInterval = readAutoInterval
			viper.Set("word.read_auto_interval", config.WordReadAutoInterval)

			config.WordSelectMode = wordSelectModeSelect.Selected
			viper.Set("word.select_mode", config.WordSelectMode)

			if showEnglishCheck.Checked != config.WordEnglishShow {
				needUpdateWord = true
			}
			config.WordEnglishShow = showEnglishCheck.Checked
			viper.Set("word.show_english", config.WordEnglishShow)

			if showChineseCheck.Checked != config.WordChineseShow {
				needUpdateWord = true
			}
			config.WordChineseShow = showChineseCheck.Checked
			viper.Set("word.show_chinese", config.WordChineseShow)

			if err := viper.WriteConfig(); err != nil {
				logrus.Errorf("failed to update config file: %v", err)
			}

			if needUpdateList {
				err := mainWindow.UpdateList()
				if err != nil {
					mainWindow.showError(err)
				}
				mainWindow.currentWord()
				needUpdateWord = true
			}
			if needUpdateWord {
				mainWindow.showWord()
			}
			if needUpdateReadButtonIcon {
				mainWindow.updateReadButtonIcon()
			}
			s.window.Close()
		},
		OnCancel: func() {
			s.window.Close()
		},
	}

	s.window.SetContent(form)
	s.window.SetIcon(resource.EarIcoe)
	return s
}
