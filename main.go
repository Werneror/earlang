package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"fyne.io/fyne/v2/app"
	"github.com/flopp/go-findfont"
	"github.com/sirupsen/logrus"
	"github.com/werneror/earlang/config"
	"github.com/werneror/earlang/window"
)

func main() {
	earlang := app.New()
	mainWindow := window.NewMainWindow(earlang)
	mainWindow.ShowAndRun()
}

func init() {
	rand.Seed(time.Now().Unix())

	if os.Getenv("FYNE_SCALE") == "" && config.FyneScale != 0 {
		err := os.Setenv("FYNE_SCALE", fmt.Sprintf("%f", config.FyneScale))
		if err != nil {
			logrus.Errorf("failed to set env FYNE_SCALE to %f: %v", config.FyneScale, err)
		}
	}

	r := os.Getenv("FYNE_FONT")
	if r != "" {
		logrus.Debugf("the value of env FYNE_FONT is %s", r)
		return
	}
	fontPath := config.FyneFont
	if fontPath == "" {
		fontPaths := findfont.List()
		for _, path := range fontPaths {
			if strings.HasSuffix(path, "HGH_CNKI.TTF") || // 光华黑体
				strings.HasSuffix(path, "simhei.ttf") || // 中易黑体
				strings.HasSuffix(path, "simkai.ttf") { // 中易楷体
				fontPath = path
				break
			}
		}
	}
	if fontPath == "" {
		logrus.Warn("the value of env FYNE_FONT is not set, and Chinese characters may not be displayed")
		return
	}
	err := os.Setenv("FYNE_FONT", fontPath)
	if err != nil {
		logrus.Errorf("failed to set env FYNE_FONT to %s: %v", fontPath, err)
	} else {
		logrus.Debugf("set the value of env FYNE_FONT to %s", fontPath)
	}
}
