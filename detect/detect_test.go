package detect

import (
	"fmt"
	"testing"
)

// func TestDetect_Compare(t *testing.T) {
// 	d := New("./../test_files/DSS_CUSTOM")
// 	d.AddOptionDir(
// 		"./../test_files/DSS_ORIGINAL_OPTION1",
// 		"./../test_files/DSS_ORIGINAL_OPTION2",
// 	)
//
// 	entries, err := d.Compare()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	fmt.Println(entries)
// }

func TestDetect_Compare(t *testing.T) {
	d := New("/Volumes/REG/Office/")
	// d.AddOptionDir(
	// 	"./../test_files/DSS_ORIGINAL_OPTION1",
	// 	"./../test_files/DSS_ORIGINAL_OPTION2",
	// )
	// err := d.AddOptionDirs("/Volumes/REG/Options/")
	err := d.AddOptionDirs("./../test_files/Options")
	if err != nil {
		t.Fatal(err)
	}
	err = d.AddSampleDirs("/Volumes/REG/Samples/")
	if err != nil {
		t.Fatal(err)
	}

	entries, err := d.Compare()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(entries)
}
