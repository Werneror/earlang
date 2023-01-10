package config

import (
	_ "embed"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/werneror/earlang/word/builtin"
)

//go:embed resources/wrong_tone.wav
var wrongToneWav []byte

var BaseDir string
var WordDir string
var PictureDir string
var WrongTonePath string

var LogLevel = "warning"
var LogFile = "earlang.log"

var PicPicker = "bing"
var PicTotalNumber = 10
var PicNumPerLine = 5
var PicMinWidth float32 = 200
var PicMinHeight float32 = 200

var PronPicker = "cambridge"
var PronRegion = "us"

const (
	WordSelectModeOrder  = "order"
	WordSelectModeRandom = "random"

	WordReadModeAuto   = "auto"
	WordReadModeOnce   = "once"
	WordReadModeManual = "manual"

	ExamineModeAll        = "all"
	ExamineModeLearned    = "learned"
	ExamineModeUnfamiliar = "unfamiliar"
)

var WordGroupName = builtin.Groups[0].Name
var WordSelectMode = WordSelectModeRandom
var WordReadMode = WordReadModeAuto
var WordReadAutoInterval = 2
var WordEnglishShow = false
var WordChineseShow = false

var WordListFileExtension = ".txt"
var WordLearnedFileExtension = ".ld"
var WordProcessFileExtension = ".proc"
var WordExamineOptionsCount = 4

var ExamineMode = ExamineModeLearned
var ExamineDataFile = "examine.json"

var UnfamiliarDataFile = "unfamiliar.json"

var FyneFont string
var FyneScale = 1.2

func init() {
	// Determine the base directory
	appData := os.Getenv("AppData")
	if _, err := os.Stat(appData); os.IsNotExist(err) {
		appData, err = os.UserHomeDir()
		if err != nil {
			appData = "."
		}
	}
	BaseDir = filepath.Join(appData, "earlang")
	if _, err := os.Stat(BaseDir); os.IsNotExist(err) {
		_ = os.MkdirAll(BaseDir, os.ModePerm)
	}
	logrus.Debugf("base directory is %s", BaseDir)

	WordDir = filepath.Join(BaseDir, "word")
	if _, err := os.Stat(WordDir); os.IsNotExist(err) {
		_ = os.MkdirAll(WordDir, os.ModePerm)
	}
	logrus.Debugf("word directory is %s", WordDir)

	err := builtin.SaveToDisk(WordDir)
	if err != nil {
		logrus.Errorf("failed to save built-in groups: %v", err)
	}

	PictureDir = filepath.Join(BaseDir, "picture")
	if _, err := os.Stat(PictureDir); os.IsNotExist(err) {
		_ = os.MkdirAll(PictureDir, os.ModePerm)
	}
	logrus.Debugf("picture directory is %s", PictureDir)

	WrongTonePath = filepath.Join(BaseDir, "wrong_tone.wav")
	if _, err := os.Stat(WrongTonePath); os.IsNotExist(err) {
		err := os.WriteFile(WrongTonePath, wrongToneWav, os.ModePerm)
		if err != nil {
			logrus.Errorf("failed to save wront tone wav to %s: %v", WrongTonePath, err)
		}
	}
	logrus.Debugf("wrong tone sound is is %s", WrongTonePath)

	// configure
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(BaseDir)

	viper.SetDefault("log.level", LogLevel)
	viper.SetDefault("log.file", LogFile)

	viper.SetDefault("picture.picker", PicPicker)
	viper.SetDefault("picture.total_number", PicTotalNumber)
	viper.SetDefault("picture.number_per_line", PicNumPerLine)
	viper.SetDefault("picture.min_height", PicMinHeight)
	viper.SetDefault("picture.min_width", PicMinWidth)

	viper.SetDefault("pronunciation.picker", PronPicker)
	viper.SetDefault("pronunciation.region", PronRegion)

	viper.SetDefault("word.group_name", WordGroupName)
	viper.SetDefault("word.select_mode", WordSelectMode)
	viper.SetDefault("word.read_mode", WordReadMode)
	viper.SetDefault("word.read_auto_interval", WordReadAutoInterval)
	viper.SetDefault("word.show_english", WordEnglishShow)
	viper.SetDefault("word.show_chinese", WordChineseShow)

	viper.SetDefault("word.examine_mode", ExamineMode)

	viper.SetDefault("fyne.font", FyneFont)
	viper.SetDefault("fyne.scale", FyneScale)

	if err := viper.ReadInConfig(); err != nil {
		logrus.Errorf("failed to read config: %v", err)
		if err := viper.SafeWriteConfig(); err != nil {
			logrus.Errorf("failed to write config: %v", err)
		}
	} else {
		// there may be new configuration items
		if err := viper.WriteConfig(); err != nil {
			logrus.Errorf("failed to update config file: %v", err)
		}
	}

	LogLevel = viper.GetString("log.level")
	LogFile = viper.GetString("log.file")

	PicPicker = viper.GetString("picture.picker")
	PicTotalNumber = viper.GetInt("picture.total_number")
	PicNumPerLine = viper.GetInt("picture.number_per_line")
	PicMinHeight = float32(viper.GetFloat64("picture.min_height"))
	PicMinWidth = float32(viper.GetFloat64("picture.min_width"))

	PronPicker = viper.GetString("pronunciation.picker")
	PronRegion = viper.GetString("pronunciation.region")

	WordGroupName = viper.GetString("word.group_name")
	WordSelectMode = viper.GetString("word.select_mode")
	WordReadMode = viper.GetString("word.read_mode")
	WordReadAutoInterval = viper.GetInt("word.read_auto_interval")
	WordEnglishShow = viper.GetBool("word.show_english")
	WordChineseShow = viper.GetBool("word.show_chinese")

	ExamineMode = viper.GetString("word.examine_mode")

	FyneFont = viper.GetString("fyne.font")
	FyneScale = viper.GetFloat64("fyne.scale")

	// log
	level, err := logrus.ParseLevel(LogLevel)
	if err != nil {
		logrus.Errorf("invalid log level %s", LogLevel)
	} else {
		logrus.SetLevel(level)
	}
	logFile := filepath.Join(BaseDir, LogFile)
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		logrus.Errorf("failed to open log file %s: %v", LogLevel, err)
	} else {
		logrus.SetOutput(f)
	}
}
