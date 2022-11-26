package main

import (
	"earlang/window"
	"os"
	"strings"

	"fyne.io/fyne/v2/app"
	"github.com/flopp/go-findfont"
)

var version = "0.0.2"

func main() {
	earlang := app.NewWithID("wiki.werner.earlang")
	mainWindow := window.NewMainWindow(earlang, version)
	mainWindow.ShowAndRun()
}

func init() {
	fontPaths := findfont.List()
	for _, path := range fontPaths {
		if strings.Contains(path, "msyh.ttc") || // 微软雅黑
			strings.Contains(path, "simkai.ttf") || // 楷体
			strings.Contains(path, "simhei.ttf") { // 黑体
			os.Setenv("FYNE_FONT", path)
			break
		}
	}
}
