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

//
