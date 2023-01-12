package window

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/werneror/earlang/resource"
	"github.com/werneror/earlang/unfamiliar"
	"github.com/werneror/earlang/word"
)

type unfamiliarWindow struct {
	window fyne.Window
}

func (u *unfamiliarWindow) showError(err error) {
	dialog.ShowError(err, u.window)
}

func (u *unfamiliarWindow) Show() {
	u.window.RequestFocus()
	u.window.CenterOnScreen()
	u.window.Show()
}

func newUnfamiliarWindow(app fyne.App, u *unfamiliar.Unfamiliar) *unfamiliarWindow {
	w := &unfamiliarWindow{
		window: app.NewWindow("EarLang Unfamiliar Words"),
	}
	w.window.SetIcon(resource.EarIcoe)

	text := ""
	for _, w := range u.AllWords() {
		text = fmt.Sprintf("%s%s,%s\n", text, w.English, w.Chinese)
	}
	textArea := widget.NewMultiLineEntry()
	textArea.SetText(text)
	textArea.SetMinRowsVisible(20)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "unfamiliar words", Widget: textArea},
		},
		OnSubmit: func() {
			newWords := make([]word.Word, 0)
			for _, line := range strings.Split(textArea.Text, "\n") {
				line = strings.TrimSpace(line)
				line = strings.Replace(line, "，", ",", 1)
				if line == "" {
					continue
				}
				pieces := strings.SplitN(line, ",", 2)
				newWord := word.Word{
					English: strings.TrimSpace(pieces[0]),
				}
				if len(pieces) > 1 {
					newWord.Chinese = strings.TrimSpace(pieces[1])
				}
				newWords = append(newWords, newWord)
			}
			err := u.Set(newWords)
			if err != nil {
				w.showError(err)
			} else {
				w.window.Close()
			}
		},
		OnCancel: func() {
			w.window.Close()
		},
	}

	w.window.SetContent(form)
	w.window.Resize(fyne.NewSize(350, 0))
	return w
}
