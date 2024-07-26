package dss

import (
	"fmt"
	"testing"
)

func TestDSS_WalkDir(t *testing.T) {
	d := New("dss", 30)
	err := d.WalkDir("/Volumes/REG/Office/")
	if err != nil {
		return
	}
	fmt.Println(len(d.Data))
}
