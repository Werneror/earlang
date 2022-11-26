package config

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var BaseDir string

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
	WordGroupTypeBuiltin = "builtin"
	WordGroupTypeCustom  = "custom"

	WordSelectModeOrder  = "order"
	WordSelectModeRandom = "random"

	WordReadModeAuto   = "auto"
	WordReadModeOnce   = "once"
	WordReadModeManual = "manual"
)

var GroupType = WordGroupTypeBuiltin
var GroupName = "tools 01"
var GroupFile = ""
var WordLearnedFile = "learned.txt"
var WordProgressFile = "progress.txt"
var WordSelectMode = WordSelectModeRandom
var WordReadMode = WordReadModeAuto
var WordReadAutoInterval = 2
var WordShow = false

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

	viper.SetDefault("word.group_type", GroupFile)
	viper.SetDefault("word.group_name", GroupName)
	viper.SetDefault("word.group_file", GroupFile)
	viper.SetDefault("word.learned_file", WordLearnedFile)
	viper.SetDefault("word.progress_file", WordProgressFile)
	viper.SetDefault("word.select_mode", WordSelectMode)
	viper.SetDefault("word.read_mode", WordReadMode)
	viper.SetDefault("word.read_auto_interval", WordReadAutoInterval)
	viper.SetDefault("word.show", WordShow)

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

	GroupFile = viper.GetString("word.group_type")
	GroupName = viper.GetString("word.group_name")
	GroupFile = viper.GetString("word.group_file")
	WordLearnedFile = viper.GetString("word.learned_file")
	WordProgressFile = viper.GetString("word.progress_file")
	WordSelectMode = viper.GetString("word.select_mode")
	WordReadMode = viper.GetString("word.read_mode")
	WordReadAutoInterval = viper.GetInt("word.read_auto_interval")
	WordShow = viper.GetBool("word.show")

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
