package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"log"
)

func (u *ui) configWindow() *fyne.Container {

	// fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
	//		layout.NewSpacer(), projectName, layout.NewSpacer())

	name := widget.NewEntry()
	name.SetPlaceHolder("Option/Sample Name")

	if u.option != nil {
		name.SetText(u.option.Name)
	}

	f := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: name},
		},
		OnSubmit:   u.saveConfig,
		OnCancel:   nil,
		SubmitText: "",
		CancelText: "",
	}

	return fyne.NewContainerWithLayout(layout.NewHBoxLayout(), f)
}

func (u *ui) saveConfig() {
	t := []string{"A", "B"}
	log.Println(t)
}
