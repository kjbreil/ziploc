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
		if u.o.InstanceDir != nil {
			items = append(items, &widget.Button{
				DisableableWidget: widget.DisableableWidget{},
				Text:              "From Instance",
				Icon:              theme.DocumentSaveIcon(),
				Importance:        widget.MediumImportance,
				OnTapped: func() {
					// get the files
					err := u.o.WalkInstance(*u.o.InstanceDir)
					if err != nil {
						dialog.ShowError(err, u.w)
						return
					}
					err = u.o.FromInstance()
					if err != nil {
						dialog.ShowError(err, u.w)
					} else {
						dialog.ShowInformation("Success", "Zip file created", u.w)
					}
				},
			})
		}
		items = append(items, &widget.Button{
			DisableableWidget: widget.DisableableWidget{},
			Text:              "Build Zip",
			Icon:              theme.DocumentSaveIcon(),
			Importance:        widget.HighImportance,
			OnTapped: func() {
				// get the files
				err := u.o.GetBaseFiles()
				if err != nil {
					dialog.ShowError(err, u.w)
					return
				}
				// make dss object for base files
				err = u.o.GetBaseDSS()
				if err != nil {
					dialog.ShowError(err, u.w)
					return
				}
				err = u.o.FromBase(filepath.Dir(u.configLocation), true)
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
