package main

import (
	"earlang/window"
	"fyne.io/fyne/v2/app"
)

var version = "0.0.2"

func main() {
	earlang := app.NewWithID("wiki.werner.earlang")
	mainWindow := window.NewMainWindow(earlang, version)
	mainWindow.ShowAndRun()
}
