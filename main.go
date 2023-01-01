package main

import (
	"os"
	"strings"

	"fyne.io/fyne/v2/app"
	"github.com/flopp/go-findfont"
	"github.com/sirupsen/logrus"
	"github.com/werneror/earlang/window"
)

func main() {
	earlang := app.New()
	mainWindow := window.NewMainWindow(earlang)
	mainWindow.ShowAndRun()
}

func init() {
	if os.Getenv("FYNE_SCALE") == "" {
		err := os.Setenv("FYNE_SCALE", "1.2")
		if err != nil {
			logrus.Errorf("failed to set env FYNE_SCALE to 1.2: %v", err)
		}
	}

	r := os.Getenv("FYNE_FONT")
	if r != "" {
		logrus.Debugf("the value of env FYNE_FONT is %s", r)
		return
	}
	fontPaths := findfont.List()
	for _, path := range fontPaths {
		if strings.HasSuffix(path, "HGH_CNKI.TTF") || // 光华黑体
			strings.HasSuffix(path, "simhei.ttf") || // 中易黑体
			strings.HasSuffix(path, "simkai.ttf") { // 中易楷体
			err := os.Setenv("FYNE_FONT", path)
			if err != nil {
				logrus.Errorf("failed to set env FYNE_FONT to %s: %v", path, err)
			} else {
				logrus.Debug("set the value of env FYNE_FONT to %s", path)
			}
			return
		}
	}
	logrus.Warn("the value of env FYNE_FONT is not set, and Chinese characters may not be displayed")
}
