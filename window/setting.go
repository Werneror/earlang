package window

import (
	"earlang/config"
	"earlang/resource"
	"earlang/word"
	"earlang/word/group"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
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

func GetGroupDisplayName(g group.Group) string {
	total := len(g.Words)
	processFilePath := filepath.Join(config.BaseDir, "word", fmt.Sprintf("%s_%s", g.Name, config.WordProgressFile))
	process, err := word.LoadPointerFromFile(processFilePath)
	if err != nil {
		logrus.Errorf("failed to load pointer from file %s: %v", processFilePath, err)
		process = 0
	}
	return fmt.Sprintf("%s (%d/%d)", g.Name, process, total)
}

func DisplayNameToGroupName(displayName string) string {
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

	groupTypeSelect := widget.NewSelect([]string{config.WordGroupTypeBuiltin, config.WordGroupTypeCustom}, func(s string) {})
	groupTypeSelect.SetSelected(config.GroupType)

	groupNames := make([]string, 0, len(group.Groups))
	var currentDisplayName string
	for _, g := range group.Groups {
		displayName := GetGroupDisplayName(g)
		groupNames = append(groupNames, displayName)
		if config.GroupName == g.Name {
			currentDisplayName = displayName
		}
	}
	groupNameSelect := widget.NewSelect(groupNames, func(s string) {})
	groupNameSelect.SetSelected(currentDisplayName)

	newGroupFile := config.GroupFile
	var groupFileButton *widget.Button
	groupFileButton = widget.NewButton(config.GroupFile, func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, s.window)
				return
			}
			if reader == nil {
				return
			}
			defer reader.Close()
			relPath, err := filepath.Rel(config.BaseDir, reader.URI().Path())
			if err != nil {
				// 如果选中的文件和 BaseDir 不在同一个磁盘上求相对路径就会出错，这时直接使用绝对路径
				newGroupFile = reader.URI().Path()
			} else {
				newGroupFile = relPath
			}
			groupFileButton.SetText(newGroupFile)
		}, s.window)
		dir, err := storage.ListerForURI(storage.NewFileURI(config.BaseDir))
		if err == nil {
			fd.SetLocation(dir)
		}
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".txt"}))
		fd.Show()
	})

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
			{Text: "custom group file", Widget: groupFileButton},
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

			if config.GroupType != groupTypeSelect.Selected {
				needUpdateList = true
			}
			config.GroupType = groupTypeSelect.Selected
			viper.Set("word.group_type", config.GroupType)

			selectedGroupName := DisplayNameToGroupName(groupNameSelect.Selected)
			if config.GroupType == config.WordGroupTypeBuiltin && config.GroupName != selectedGroupName {
				needUpdateList = true
			}
			config.GroupName = selectedGroupName
			viper.Set("word.group_name", config.GroupName)

			if config.GroupType == config.WordGroupTypeCustom && config.GroupFile != newGroupFile {
				needUpdateList = true
			}
			config.GroupFile = newGroupFile
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
