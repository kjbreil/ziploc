package iniUpdater

import (
	"fmt"
	"testing"
)

func TestOption_ReadIni(t *testing.T) {
	// /Users/kjell/working/SMSLane/Options/Application_SMSPro/FirstLoad/Office
	// u, err := CompareTwo("/Users/kjell/working/SMSLane/Office/Setting.ini", "../../ncbp_lab/lane/OFFICE/Setting.ini")
	u, err := CompareTwo("/Users/kjell/working/SMSLane/Options/Application_SMSPro/FirstLoad/Office/System.ini", "/Users/kjell/working/SMSLane/Office/System.ini")
	if err != nil {
		t.Fatal(err)
	}

	b := u.Bytes()

	fmt.Println(string(b))
	t.Fail()
}
