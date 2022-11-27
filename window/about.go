package window

import (
	"earlang/resource"
	"image/color"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type aboutWindow struct {
	window fyne.Window
}

func (a *aboutWindow) Show() {
	a.window.RequestFocus()
	a.window.Show()
}

func newAboutWindow(app fyne.App) *aboutWindow {
	w := &aboutWindow{
		window: app.NewWindow("EarLang About"),
	}

	softwareText := canvas.NewText("EarLang", color.Black)
	softwareText.TextSize = 24
	versionText := canvas.NewText(app.Metadata().Version, color.Black)
	softwareBox := container.New(layout.NewHBoxLayout(), softwareText, versionText)

	homepageText := canvas.NewText("Homepage:", color.Black)
	githubUrl := "https://github.com/werneror/earlang"
	homepageUrl, _ := url.Parse(githubUrl)
	homepageLink := widget.NewHyperlink(githubUrl, homepageUrl)
	homepageBox := container.New(layout.NewHBoxLayout(), homepageText, homepageLink)

	AuthorText := canvas.NewText("Author:   werner <me@werner.wiki>", color.Black)

	AboutBox := container.New(layout.NewVBoxLayout(), softwareBox, homepageBox, AuthorText)

	w.window.SetContent(container.New(layout.NewCenterLayout(), AboutBox))
	w.window.SetIcon(resource.EarIcoe)
	w.window.Resize(fyne.NewSize(459, 165))
	return w
}
