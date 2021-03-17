package gui

import (
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
	"strconv"
)

func (u *ui) mainTab() *fyne.Container {
	name := &widget.Entry{Text: u.o.Name}

	priority := &widget.Entry{
		Text: strconv.Itoa(u.o.Priority),
		Validator: func(s string) error {
			i, err := strconv.Atoi(s)
			u.o.Priority = i
			return err
		},
	}

	t := "Option"
	if !u.o.IsOption {
		t = "Sample"
	}
	kind := &widget.Select{
		DisableableWidget: widget.DisableableWidget{},
		Selected:          t,
		Options:           []string{"Option", "Sample"},
	}

	baseFolder := &widget.Entry{Text: u.o.BaseFolder}
	outFolder := &widget.Entry{Text: u.o.ZipOut}

	items := []*widget.FormItem{
		{Text: "Name", Widget: name},
		{Text: "Kind", Widget: kind},
		{Text: "Priority", Widget: priority},
		{Text: "Source Folder", Widget: baseFolder},
		{Text: "Publish Folder", Widget: outFolder},
		// {Text: "Instance?", u.instanceSelect()},
	}
	items = append(items, &widget.FormItem{Text: "Instance Name", Widget: u.instanceSelect()})

	instanceName := &widget.Entry{}
	if u.o.Instance != nil {
		instanceName.SetText(*u.o.Instance)
		items = append(items, &widget.FormItem{Text: "Instance Name", Widget: instanceName})
	}

	instanceDir := &widget.Entry{}
	if u.o.InstanceDir != nil {
		instanceDir.SetText(*u.o.InstanceDir)
		items = append(items, &widget.FormItem{Text: "Instance Dir", Widget: instanceDir})
	}

	f := &widget.Form{
		Items: items,
		OnSubmit: func() {
			u.o.Name = name.Text
			u.o.BaseFolder = baseFolder.Text
			u.o.ZipOut = outFolder.Text
			in, id := instanceName.Text, instanceDir.Text
			u.o.Instance = &in
			u.o.InstanceDir = &id
			if kind.Selected == "Sample" {
				u.o.IsOption = false
			} else {
				u.o.IsOption = true
			}

			u.o.WriteConfig(u.configLocation)
			u.configWindow()

			dialog.ShowInformation("Success", "Config file saved", u.w)
		},
		SubmitText: "Save",
		OnCancel: func() {
			u.o.Name = name.Text
			u.o.BaseFolder = baseFolder.Text
			u.o.ZipOut = outFolder.Text
			in, id := instanceName.Text, instanceDir.Text
			u.o.Instance = &in
			u.o.InstanceDir = &id
			if kind.Selected == "Sample" {
				u.o.IsOption = false
			} else {
				u.o.IsOption = true
			}
			dialog.ShowInformation("Success", "Config was set in memory", u.w)
		},
		CancelText: "Set",
	}

	return container.NewCenter(f)
}

func (u *ui) instanceSelect() fyne.CanvasObject {
	fromInstance := widget.NewCheck("Instance", func(b bool) {
		if !b {
			dialog.ShowConfirm("Instance Reset", "Are you sure you want to clear instence info", func(b bool) {
				if b {
					u.o.Instance = nil
					u.o.InstanceDir = nil
				}
			}, u.w)
		} else {
			u.o.DefaultExclude()
			u.o.DefaultInclude()
			u.o.Instance = new(string)
			u.o.InstanceDir = new(string)
		}
		u.configWindow()
	})
	fromInstance.Checked = u.o.Exclude != nil
	return fromInstance
}
