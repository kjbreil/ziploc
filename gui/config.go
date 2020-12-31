package gui

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"log"
	"strconv"
)

func (u *ui) configWindow() {
	if u.o == nil {
		log.Panic("config window called without an o loaded")
	}

	name := &widget.Entry{Text: u.o.Name}

	priority := &widget.Entry{
		DisableableWidget: widget.DisableableWidget{},
		Text:              strconv.Itoa(u.o.Priority),
		Validator: func(s string) error {
			i, err := strconv.Atoi(s)
			u.o.Priority = i
			return err
		},
	}

	baseFolder := &widget.Entry{Text: u.o.BaseFolder}
	f := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: name},
			{Text: "Priority", Widget: priority},
			{Text: "Source Folder", Widget: baseFolder},
		},
		OnSubmit: func() {
			u.o.Name = name.Text
			u.o.BaseFolder = baseFolder.Text

			u.o.WriteConfig(u.configLocation)
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
	for _, ee := range u.o.Exclude {
		exclude := &widget.Entry{Text: ee}
		excluded = append(excluded, exclude)
	}
	return fyne.NewContainerWithLayout(layout.NewVBoxLayout(), excluded...)
}
