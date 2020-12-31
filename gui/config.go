package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"log"
	"strconv"
)

func (u *ui) configWindow() {
	if u.option == nil {
		log.Panic("config window called without an option loaded")
	}

	name := &widget.Entry{Text: u.option.Name}

	priority := &widget.Entry{
		DisableableWidget: widget.DisableableWidget{},
		Text:              strconv.Itoa(u.option.Priority),
		Validator: func(s string) error {
			i, err := strconv.Atoi(s)
			u.option.Priority = i
			return err
		},
	}

	baseFolder := &widget.Entry{Text: u.option.BaseFolder}
	f := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: name},
			{Text: "Priority", Widget: priority},
			{Text: "Source Folder", Widget: baseFolder},
		},
		OnSubmit: func() {
			u.option.Name = name.Text
			u.option.BaseFolder = baseFolder.Text

			u.option.WriteConfig(u.configLocation)
			u.configWindow()
		},
		OnCancel: nil,
	}

	var con []fyne.CanvasObject
	con = append(con, f)

	con = append(con, u.showExcludes())

	u.loadWindow(fyne.NewContainerWithLayout(layout.NewHBoxLayout(), con...))
}

func (u *ui) showExcludes() *fyne.Container {
	var excluded []fyne.CanvasObject
	for _, ee := range u.option.Exclude {
		exclude := &widget.Entry{Text: ee}
		excluded = append(excluded, exclude)
	}
	return fyne.NewContainerWithLayout(layout.NewVBoxLayout(), excluded...)
}
