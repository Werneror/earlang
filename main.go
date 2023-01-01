package main

import (
	"os"
	"strings"

	"fyne.io/fyne/v2/app"
	"github.com/flopp/go-findfont"
	"github.com/werneror/earlang/window"
)

func main() {
	earlang := app.New()
	mainWindow := window.NewMainWindow(earlang)
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
