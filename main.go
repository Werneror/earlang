package main

import (
	"earlang/window"
	"fyne.io/fyne/v2/app"
)

var version = "0.0.1"

func main() {
	earlang := app.New()
	mainWindow := window.NewMainWindow(earlang, version)
	mainWindow.ShowAndRun()
}
