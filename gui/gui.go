package gui

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"github.com/kjbreil/ziploc/option"
	"log"
)

type ui struct {
	o *option.Option

	a fyne.App
	w fyne.Window

	configLocation string
	// toolbar *fyne.Container
}

func OpenGui() {
	u := new(ui)
	u.init()

	u.a = app.New()
	u.w = u.a.NewWindow("Box Layout")

	u.loadWindow()
	u.w.Resize(fyne.NewSize(640, 460))

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
		log.Panic(err)
	}
	u.configWindow()

}

func (u *ui) loadWindow(containers ...fyne.CanvasObject) {
	var newCont []fyne.CanvasObject
	newCont = append(newCont, u.createToolbar())
	if len(containers) == 0 {
		newCont = append(newCont, layout.NewSpacer())
	}
	newCont = append(newCont, containers...)
	newCont = append(newCont, u.createFooter())

	u.w.SetContent(fyne.NewContainerWithLayout(layout.NewVBoxLayout(), newCont...))

}

// func (u *ui)
