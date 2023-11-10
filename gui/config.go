package gui

import (
	"fyne.io/fyne/v2/container"
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
	if u.o.Instance == nil {
		u.o.Instance = new(string)
	}
	if u.o.InstanceDir == nil {
		u.o.Instance = new(string)
	}
	if u.o.Exclude == nil {
		u.o.DefaultExclude()
	}
	if u.o.Include == nil {
		u.o.DefaultInclude()
	}
}
