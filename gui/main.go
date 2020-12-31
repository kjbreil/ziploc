package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/sqweek/dialog"
	"log"
)

type ziploc struct {
	name string

	update chan struct{}
}

func main() {
	z := new(ziploc)
	z.update = make(chan struct{})

	a := app.New()
	w := a.NewWindow("Box Layout")

	container := z.toolbar()

	projectName := widget.NewLabel("XXX")

	// text4 := canvas.NewText("centered", color.White)

	centered := fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
		layout.NewSpacer(), projectName, layout.NewSpacer())
	w.SetContent(fyne.NewContainerWithLayout(layout.NewVBoxLayout(), container, centered))

	w.Resize(fyne.NewSize(640, 460))

	go func() {
		// wait for update signal
		<-z.update
		projectName.SetText(z.name)
	}()

	w.ShowAndRun()
}

func (z *ziploc) toolbar() *fyne.Container {
	open := widget.NewButton("Open", func() {
		filename, _ := dialog.File().Filter("Json", "json").Load()

		z.loadJsonConfig(filename)
		// dialog.
		//
		// dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		// 	if err == nil && reader == nil {
		// 		return
		// 	}
		// 	if err != nil {
		// 		dialog.ShowError(err, win)
		// 		return
		// 	}
		// 	z.loadJsonConfig()
		// }, win)

		// fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		// 	if err == nil && reader == nil {
		// 		return
		// 	}
		// 	if err != nil {
		// 		dialog.ShowError(err, win)
		// 		return
		// 	}
		// 	z.loadJsonConfig()
		//
		// }, win)
		// fd.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".txt"}))
		// fd.Resize(fyne.NewSize(640, 460))
		//
		// fd.Show()

	})
	container := fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
		open)

	return container
}

func (z *ziploc) loadJsonConfig(filename string) {
	log.Println(filename)
	// func (z *ziploc) loadJsonConfig(f fyne.URIReadCloser) {
	z.name = filename
	z.update <- struct{}{}
}
