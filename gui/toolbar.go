package gui

import (
	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	sysDiag "github.com/sqweek/dialog"
	"path/filepath"
)

func (u *ui) createToolbar() *fyne.Container {

	var items []fyne.CanvasObject

	// Open button
	items = append(items, widget.NewButton("Open", func() {
		filename, err := sysDiag.File().Filter("Json", "json").Load()
		if err != nil {
			dialog.ShowError(err, u.w)
		} else {
			u.loadJsonConfig(&filename)
		}
	}))

	// a config is loaded so show buttons for interacting with the o
	if u.o != nil {
		// open config form button
		items = append(items, widget.NewButton("Config", func() {
			u.configWindow()
		}))
		items = append(items, widget.NewButton("Build", func() {
			err := u.o.FromBase(filepath.Dir(u.configLocation))
			if err != nil {
				dialog.ShowError(err, u.w)
			} else {
				dialog.ShowInformation("Success", "Zip file created", u.w)
			}
		}))
		// Spacer so o name is on right
		items = append(items, layout.NewSpacer())
		// Option Name
		items = append(items, widget.NewLabel(u.o.Name))
	}

	container := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), items...)

	return container
}
