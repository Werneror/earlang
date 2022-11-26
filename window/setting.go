package window

import (
	"earlang/config"
	"earlang/resource"
	"earlang/word/group"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type settingWindow struct {
	window     fyne.Window
	mainWindow *MainWindow
}

func (s *settingWindow) Show() {
	s.window.RequestFocus()
	s.window.Show()
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

	groupTypeSelect := widget.NewSelect([]string{config.WordGroupTypeBuiltin, config.WordGroupTypeCustom}, func(s string) {})
	groupTypeSelect.SetSelected(config.GroupType)

	groupNames := make([]string, 0, len(group.Groups))
	for _, g := range group.Groups {
		groupNames = append(groupNames, g.Name)
	}
	groupNameSelect := widget.NewSelect(groupNames, func(s string) {})
	groupNameSelect.SetSelected(config.GroupName)

	// TODO: use dialog.FileDialog
	groupFileEntry := widget.NewEntry()
	groupFileEntry.SetText(config.GroupFile)

	readModeSelect := widget.NewSelect([]string{config.WordReadModeAuto, config.WordReadModeOnce, config.WordReadModeManual}, func(string) {})
	readModeSelect.SetSelected(config.WordReadMode)

	wordSelectModeSelect := widget.NewSelect([]string{config.WordSelectModeOrder, config.WordSelectModeRandom}, func(string) {})
	wordSelectModeSelect.SetSelected(config.WordSelectMode)

	readAutoIntervalSelect := widget.NewSelect([]string{"1", "2", "3", "4", "5", "10"}, func(string) {})
	readAutoIntervalSelect.SetSelected(fmt.Sprintf("%d", config.WordReadAutoInterval))

	showEnglishCheck := widget.NewCheck("show word below the pictures", func(bool) {})
	showEnglishCheck.SetChecked(config.WordEnglishShow)

	showChineseCheck := widget.NewCheck("show Chinese below the pictures", func(bool) {})
	showChineseCheck.SetChecked(config.WordChineseShow)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "log level", Widget: logLevelSelect},
			{Text: "pronunciation source", Widget: pronPickerSelect},
			{Text: "pronunciation region", Widget: pronRegionSelect},
			{Text: "picture source", Widget: picPickerSelect},
			{Text: "picture number", Widget: picNumberSelect},
			{Text: "word group type", Widget: groupTypeSelect},
			{Text: "builtin group name", Widget: groupNameSelect},
			{Text: "custom group file", Widget: groupFileEntry},
			{Text: "word read mode", Widget: readModeSelect},
			{Text: "word auto read interval(s)", Widget: readAutoIntervalSelect},
			{Text: "word select mode", Widget: wordSelectModeSelect},
			{Text: "show word", Widget: showEnglishCheck},
			{Text: "show Chinese", Widget: showChineseCheck},
		},
		OnSubmit: func() {
			needUpdateList := true
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

			if config.GroupType != groupTypeSelect.Selected {
				needUpdateList = true
			}
			config.GroupType = groupTypeSelect.Selected
			viper.Set("word.group_type", config.GroupType)

			if config.GroupType == config.WordGroupTypeBuiltin && config.GroupName != groupNameSelect.Selected {
				needUpdateList = true
			}
			config.GroupName = groupNameSelect.Selected
			viper.Set("word.group_name", config.GroupName)

			if config.GroupType == config.WordGroupTypeCustom && config.GroupFile != groupFileEntry.Text {
				needUpdateList = true
			}
			config.GroupFile = groupFileEntry.Text
			viper.Set("word.group_file", config.GroupFile)

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
