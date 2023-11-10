package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	sysDiag "github.com/sqweek/dialog"
)

func (u *ui) createToolbar() *fyne.Container {

	var items []fyne.CanvasObject

	// Open button
	items = append(items, &widget.Button{
		DisableableWidget: widget.DisableableWidget{},
		Text:              "Open",
		Icon:              theme.FileIcon(),
		Importance:        widget.MediumImportance,
		OnTapped: func() {
			filename, err := sysDiag.File().Filter("Json", "json").Load()
			if err != nil {
				dialog.ShowError(err, u.w)
			} else {
				u.loadJsonConfig(&filename)
			}
		},
	})

	// a config is loaded so show buttons for interacting with the o
	if u.o != nil {
		// open config form button
		items = append(items, &widget.Button{
			DisableableWidget: widget.DisableableWidget{},
			Text:              "Config",
			Icon:              theme.SettingsIcon(),
			Importance:        widget.LowImportance,
			OnTapped: func() {
				u.configWindow()
			},
		})
		// Spacer so option name is on right
		items = append(items, layout.NewSpacer())
		// Option Name
		items = append(items, widget.NewLabel(u.o.Name))
	}

	return fyne.NewContainerWithLayout(layout.NewHBoxLayout(), items...)
}
