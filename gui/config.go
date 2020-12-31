package gui

import (
	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
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

	t := "Option"
	if u.o.IsSample {
		t = "Sample"
	}
	kind := &widget.Select{
		DisableableWidget: widget.DisableableWidget{},
		Selected:          t,
		Options:           []string{"Option", "Sample"},
	}

	baseFolder := &widget.Entry{Text: u.o.BaseFolder}
	outFolder := &widget.Entry{Text: u.o.ZipOut}
	f := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: name},
			{Text: "Kind", Widget: kind},
			{Text: "Priority", Widget: priority},
			{Text: "Source Folder", Widget: baseFolder},
			{Text: "Publish Folder", Widget: outFolder},
		},
		OnSubmit: func() {
			u.o.Name = name.Text
			u.o.BaseFolder = baseFolder.Text
			u.o.ZipOut = outFolder.Text
			if kind.Selected == "Sample" {
				u.o.IsSample = true
			} else {
				u.o.IsSample = false
			}

			u.o.WriteConfig(u.configLocation)
			u.configWindow()

			dialog.ShowInformation("Success", "Config file saved", u.w)
		},
		OnCancel: nil,
	}

	var con []fyne.CanvasObject

	con = append(con, f)
	var accItems []*widget.AccordionItem
	if u.o.Exclude != nil {
		accItems = append(accItems, widget.NewAccordionItem("Excludes", u.showExcludes()))
	}

	if u.o.Include != nil {
		accItems = append(accItems, widget.NewAccordionItem("Includes", u.showIncludes()))
	}

	if len(accItems) > 0 {
		acc := widget.NewAccordion(accItems...)
		con = append(con, acc)
	}

	who := fyne.NewContainerWithLayout(layout.NewVBoxLayout(),
		u.optionOptions(),
		fyne.NewContainerWithLayout(layout.NewHBoxLayout(), con...),
	)

	u.loadWindow(who)
}

func (u *ui) optionOptions() *fyne.Container {
	fromInstance := widget.NewCheck("Instance", func(b bool) {
		if !b {
			u.o.Exclude = nil
			u.o.Include = nil
		} else {
			u.o.DefaultExclude()
			u.o.DefaultInclude()
		}
		u.configWindow()
	})
	fromInstance.Checked = u.o.Exclude != nil
	return fyne.NewContainerWithLayout(layout.NewHBoxLayout(), fromInstance)
}

func (u *ui) showExcludes() *fyne.Container {
	var excluded []fyne.CanvasObject
	for _, ee := range *u.o.Exclude {
		exclude := &widget.Entry{Text: ee}
		excluded = append(excluded, exclude)
	}
	return fyne.NewContainerWithLayout(layout.NewVBoxLayout(), excluded...)
}
func (u *ui) showIncludes() *fyne.Container {
	var included []*widget.AccordionItem
	for section, eis := range u.o.Include {
		var sectionItems []fyne.CanvasObject
		for _, ei := range eis {
			include := &widget.Entry{Text: ei}
			sectionItems = append(sectionItems, include)
		}

		con := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), sectionItems...)

		included = append(included, widget.NewAccordionItem(section, con))
	}

	acc := widget.NewAccordion(included...)

	return fyne.NewContainerWithLayout(layout.NewVBoxLayout(), acc)
}
