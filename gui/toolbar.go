package gui

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/sqweek/dialog"
	"log"
)

func (u *ui) createToolbar() *fyne.Container {

	var items []fyne.CanvasObject

	// Open button
	items = append(items, widget.NewButton("Open", func() {
		filename, err := dialog.File().Filter("Json", "json").Load()
		if err != nil {
			log.Panic(err)
		}
		u.loadJsonConfig(&filename)
	}))

	// a config is loaded so show buttons for interacting with the o
	if u.o != nil {
		// open config form button
		items = append(items, widget.NewButton("Config", func() {
			u.configWindow()
		}))
		// Spacer so o name is on right
		items = append(items, layout.NewSpacer())
		// Option Name
		items = append(items, widget.NewLabel(u.o.Name))
	}

	container := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), items...)

	return container
}
