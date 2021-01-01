package gui

import (
	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"path/filepath"
)

func (u *ui) createFooter() *fyne.Container {
	var items []fyne.CanvasObject
	items = append(items, layout.NewSpacer())

	// a config is loaded so show buttons for interacting with the option
	if u.o != nil {

		items = append(items, &widget.Button{
			DisableableWidget: widget.DisableableWidget{},
			Text:              "Build Zip",
			Icon:              theme.DocumentSaveIcon(),
			Importance:        widget.HighImportance,
			Alignment:         0,
			IconPlacement:     0,
			OnTapped: func() {
				err := u.o.FromBase(filepath.Dir(u.configLocation))
				if err != nil {
					dialog.ShowError(err, u.w)
				} else {
					dialog.ShowInformation("Success", "Zip file created", u.w)
				}
			},
		})
	}

	return fyne.NewContainerWithLayout(layout.NewHBoxLayout(), items...)
}
