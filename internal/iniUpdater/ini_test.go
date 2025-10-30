package iniUpdater

import (
	"fmt"
	"testing"
)

func TestOption_ReadIni(t *testing.T) {
	u, err := CompareTwo("APIDATACAP_ORI.INI", "APIDATACAP_CUS.INI", "XCHDEV")
	if err != nil {
		t.Fatal(err)
	}

	b := u.Bytes()
	fmt.Println(string(b))

	// f, err := File("APIDATACAP_ORI.INI", "APIDATACAP_CUS.INI", "XCHDEV")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// s, err := f.Stat()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	//
	// b := make([]byte, s.Size())
	// i, err := f.Read(b)
	// if err != nil {
	// 	t.Fatal(i, err)
	// }
	// fmt.Println(string(b))

	t.Fail()
}

func TestOption_CreateManual(t *testing.T) {
	u := CreateManual("NCBP_VERSION.INI", "OFFICE", "VERSION", "TEST", "v.1.0.0.0")
	// u := New()
	// u.file = "NCBP_VERSION.INI"
	// u.folder = "OFFICE"
	// u.sections = append(u.sections, "VERSION")
	// u.sectionKeys["VERSION"] = []string{"TEST"}
	// //u.checkFor["VERSION"] = make(map[string]string)
	// //u.checkFor["VERSION"]["TEST"] = "v1.0.0.0"
	// u.D["VERSION"] = make(map[string]string)
	// u.D["VERSION"]["TEST"] = "v1.0.0.0"

	b := u.Bytes()
	fmt.Println(string(b))
}

//
