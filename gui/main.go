package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"github.com/kjbreil/ziploc"
	"log"
)

type ui struct {
	option *ziploc.Option

	a fyne.App
	w fyne.Window

	// toolbar *fyne.Container
}

func main() {
	u := new(ui)
	u.init()

	u.a = app.New()
	u.w = u.a.NewWindow("Box Layout")

	// u.toolbar = u.createToolbar()
	// u.w.SetContent(fyne.NewContainerWithLayout(layout.NewVBoxLayout(), u.toolbar))
	u.loadWindow()
	u.w.Resize(fyne.NewSize(640, 460))

	u.w.ShowAndRun()
}

func (u *ui) init() {
}

func (u *ui) loadJsonConfig(filename *string) {
	var err error

	u.option, err = ziploc.ReadConfig(filename)
	if err != nil {
		// TODO: Make this present an error dialog
		log.Panic(err)
	}
	u.loadWindow()

}

func (u *ui) loadWindow(containers ...fyne.CanvasObject) {
	containers = append([]fyne.CanvasObject{u.createToolbar()}, containers...)
	u.w.SetContent(fyne.NewContainerWithLayout(layout.NewVBoxLayout(), containers...))
}
