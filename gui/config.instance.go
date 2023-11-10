package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func (u *ui) instanceTab() fyne.CanvasObject {

	var items []*widget.FormItem
	u.initInstance()

	// TODO: Should be a dropdown
	items = append(items, &widget.FormItem{
		Text: "Instance",
		Widget: &widget.Entry{
			Text: *u.o.Instance,
		},
	})
	// TODO: Should be file picker
	items = append(items, &widget.FormItem{
		Text: "InstanceDir",
		Widget: &widget.Entry{
			Text: *u.o.InstanceDir,
		},
	})

	f := &widget.Form{
		BaseWidget: widget.BaseWidget{},
		Items:      items,
		OnSubmit:   nil,
		OnCancel:   nil,
		SubmitText: "",
		CancelText: "",
	}

	ExcInc := container.NewAppTabs(
		container.NewTabItem("Includes", u.showIncludes()),
		container.NewTabItem("Excludes", u.showExcludes()),
	)

	return container.NewBorder(f, nil, nil, nil, ExcInc)
}

func (u *ui) showExcludes() fyne.CanvasObject {
	icon := widget.NewIcon(nil)
	label := widget.NewLabel("Select An Item From The List")
	hbox := container.New(layout.NewHBoxLayout(), icon, label)

	list := widget.NewList(func() int {
		return len(u.o.Exclude)
	}, func() fyne.CanvasObject {
		return widget.NewLabel("Text")
	}, func(id widget.ListItemID, item fyne.CanvasObject) {
		item.(*widget.Label).SetText(u.o.Exclude[id])
	})
	list.OnSelected = func(id widget.ListItemID) {
		label.SetText(u.o.Exclude[id])
		icon.SetResource(theme.DocumentIcon())
	}
	list.OnUnselected = func(id widget.ListItemID) {
		label.SetText("Select An Item From The List")
		icon.SetResource(nil)
	}

	s := container.NewHSplit(list, hbox)
	s.SetOffset(0.3)
	return s
}

func (u *ui) showIncludeSection(name string) fyne.CanvasObject {

	list := widget.NewList(func() int {
		return len(u.o.Include[name])
	}, func() fyne.CanvasObject {

		return widget.NewLabel("Text")
	}, func(id widget.ListItemID, item fyne.CanvasObject) {
		item.(*widget.Label).SetText(u.o.Include[name][id])
	})
	return container.NewMax(list)
}
func (u *ui) showIncludes() fyne.CanvasObject {
	urLabel := widget.NewLabel("Select An Item From Exclude Section")
	lrLabel := widget.NewLabel("Select An Item From The List")
	ur := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), urLabel)
	lr := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), lrLabel)

	var sections []string
	for k := range u.o.Include {
		sections = append(sections, k)
	}

	list := widget.NewList(func() int {
		return len(sections)
	}, func() fyne.CanvasObject {

		return widget.NewLabel("Text")
	}, func(id widget.ListItemID, item fyne.CanvasObject) {
		item.(*widget.Label).SetText(sections[id])
	})
	list.OnSelected = func(id widget.ListItemID) {
		ur.Objects = []fyne.CanvasObject{u.showIncludeSection(sections[id])}
		ur.Refresh()
		lrLabel.SetText(sections[id])
	}
	list.OnUnselected = func(id widget.ListItemID) {
		urLabel.SetText("Done")
		lrLabel.SetText("Done")
	}
	v := container.NewVSplit(ur, lr)
	v.SetOffset(0.4)
	s := container.NewHSplit(list, v)
	s.SetOffset(0.3)

	return s
}
