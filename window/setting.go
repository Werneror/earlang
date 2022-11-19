package window

import (
	"earlang/config"
	"earlang/resource"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type settingWindow struct {
	window fyne.Window
}

func (s *settingWindow) Show() {
	s.window.RequestFocus()
	s.window.Show()
}

func newSettingWindow(app fyne.App) *settingWindow {
	s := &settingWindow{
		window: app.NewWindow("EarLang Setting"),
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

	readModeSelect := widget.NewSelect([]string{config.WordReadModeAuto, config.WordReadModeOnce, config.WordReadModeManual}, func(string) {})
	readModeSelect.SetSelected(config.WordReadMode)

	wordSelectModeSelect := widget.NewSelect([]string{config.WordSelectModeOrder, config.WordSelectModeRandom}, func(string) {})
	wordSelectModeSelect.SetSelected(config.WordSelectMode)

	readAutoIntervalSelect := widget.NewSelect([]string{"1", "2", "3", "4", "5", "10"}, func(string) {})
	readAutoIntervalSelect.SetSelected(fmt.Sprintf("%d", config.WordReadAutoInterval))

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "log level", Widget: logLevelSelect},
			{Text: "pronunciation source", Widget: pronPickerSelect},
			{Text: "pronunciation region", Widget: pronRegionSelect},
			{Text: "picture source", Widget: picPickerSelect},
			{Text: "picture number", Widget: picNumberSelect},
			{Text: "word read mode", Widget: readModeSelect},
			{Text: "word auto read interval(s)", Widget: readAutoIntervalSelect},
			{Text: "word select mode", Widget: wordSelectModeSelect},
		},
		OnSubmit: func() {
			picNumber, err := strconv.Atoi(picNumberSelect.Selected)
			if err != nil {
				logrus.Error("failed to parse picture total number %s: %v", picNumberSelect.Selected, err)
				picNumber = config.PicTotalNumber
			}
			readAutoInterval, err := strconv.Atoi(readAutoIntervalSelect.Selected)
			if err != nil {
				logrus.Error("failed to parse word read auto interval %s: %v", readAutoIntervalSelect.Selected, err)
				readAutoInterval = config.WordReadAutoInterval
			}
			config.LogLevel = logLevelSelect.Selected
			viper.Set("log.level", config.LogLevel)
			config.PronPicker = pronPickerSelect.Selected
			viper.Set("pronunciation.picker", config.PronPicker)
			config.PronRegion = pronRegionSelect.Selected
			viper.Set("pronunciation.region", config.PronRegion)
			config.PicPicker = picPickerSelect.Selected
			viper.Set("picture.picker", config.PicPicker)
			config.PicTotalNumber = picNumber
			viper.Set("picture.total_number", config.PicTotalNumber)
			config.WordReadMode = readModeSelect.Selected
			viper.Set("word.read_mode", config.WordReadMode)
			config.WordReadAutoInterval = readAutoInterval
			viper.Set("word.read_auto_interval", config.WordReadAutoInterval)
			config.WordSelectMode = wordSelectModeSelect.Selected
			viper.Set("word.select_mode", config.WordSelectMode)
			if err := viper.WriteConfig(); err != nil {
				logrus.Errorf("failed to update config file: %v", err)
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
