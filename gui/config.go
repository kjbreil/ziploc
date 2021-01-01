package gui

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"log"
)

func (u *ui) configWindow() {
	if u.o == nil {
		log.Panic("config window called without an o loaded")
	}

	who := container.NewAppTabs(
		container.NewTabItem("Main", u.mainTab()),
		container.NewTabItem("Instance", u.instanceTab()),
	)

	u.loadWindow(who)

}

func (u *ui) initInstance() {
	if u.o.Instance != nil {
		u.o.Instance = new(string)
	}
	if u.o.InstanceDir != nil {
		u.o.Instance = new(string)
	}
	if u.o.Exclude != nil {
		u.o.DefaultExclude()
	}
	if u.o.Include != nil {
		u.o.DefaultInclude()
	}
}

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
	// items = append(items, &widget.FormItem{
	// 	Text:   "InstanceDir",
	// 	Widget: u.showExcludes(),
	// })

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
	hbox := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), icon, label)

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
		log.Println(len(u.o.Include[name]))
		return len(u.o.Include[name])
	}, func() fyne.CanvasObject {

		return widget.NewLabel("Text")
	}, func(id widget.ListItemID, item fyne.CanvasObject) {
		fmt.Println(u.o.Include[name][id])
		item.(*widget.Label).SetText(u.o.Include[name][id])
	})
	// list.OnSelected = func(id widget.ListItemID) {
	// 	label.SetText(u.o.Exclude[id])
	// }
	// list.OnUnselected = func(id widget.ListItemID) {
	// 	label.SetText("Select An Item From The List")
	// }
	//
	// s := container.NewHSplit(list, hbox)
	// s.SetOffset(0.3)
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
		// ur.Objects = []fyne.CanvasObject{u.showIncludeSection(sections[id], ur, lrLabel)}
		// ur.Refresh()
		ur.Objects = []fyne.CanvasObject{u.showIncludeSection(sections[id])}
		ur.Refresh()
		lrLabel.SetText(sections[id])
	}
	list.OnUnselected = func(id widget.ListItemID) {
		urLabel.SetText("Done")
		lrLabel.SetText("Done")
	}
	//
	// var included []*widget.AccordionItem
	// for section, eis := range u.o.Include {
	// 	var sectionItems []fyne.CanvasObject
	// 	for _, ei := range eis {
	// 		include := &widget.Entry{Text: ei}
	// 		sectionItems = append(sectionItems, include)
	// 	}
	// 	// con := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), sectionItems...)
	// 	included = append(included, widget.NewAccordionItem(section, u.showIncludeSection(section, hbox, label, icon)))
	// }
	//
	// acc := widget.NewAccordion(included...)

	s := container.NewHSplit(list, container.NewVSplit(ur, lr))
	s.SetOffset(0.3)

	return s
}

func makeListTab() fyne.CanvasObject {
	var data []string
	for i := 0; i < 1000; i++ {
		data = append(data, fmt.Sprintf("Test Item %d", i))
	}

	icon := widget.NewIcon(nil)
	label := widget.NewLabel("Select An Item From The List")
	hbox := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), icon, label)

	list := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return fyne.NewContainerWithLayout(layout.NewHBoxLayout(), widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(data[id])
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		label.SetText(data[id])
		icon.SetResource(theme.DocumentIcon())
	}
	list.OnUnselected = func(id widget.ListItemID) {
		label.SetText("Select An Item From The List")
		icon.SetResource(nil)
	}
	list.Select(125)
	return container.NewHBox(list, fyne.NewContainerWithLayout(layout.NewCenterLayout(), hbox))
}
