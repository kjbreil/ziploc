package gui

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/theme"
	"github.com/kjbreil/ziploc/option"
)

type ui struct {
	o *option.Option

	a fyne.App
	w fyne.Window

	configLocation string
}

func OpenGui() {
	u := new(ui)
	u.init()

	u.a = app.New()
	u.w = u.a.NewWindow("ZipLoc")

	u.loadWindow()
	u.w.Resize(fyne.NewSize(640, 480))

	u.w.ShowAndRun()
}

func (u *ui) init() {
}

func (u *ui) loadJsonConfig(filename *string) {
	var err error

	u.configLocation = *filename

	u.o, err = option.ReadConfig(filename)
	if err != nil {
		// TODO: Make this present an error dialog
		dialog.ShowError(err, u.w)
	}

	u.configWindow()

}

func (u *ui) loadWindow(containers ...fyne.CanvasObject) {

	con := container.NewBorder(u.createToolbar(), u.createFooter(), nil, nil, container.NewMax(containers...))
	u.w.SetContent(con)

}

func makeCell() fyne.CanvasObject {
	rect := canvas.NewRectangle(theme.PrimaryColor())
	rect.SetMinSize(fyne.NewSize(5, 5))
	return rect
}

// func (u *ui)
