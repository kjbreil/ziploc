package shares

import (
	"fmt"
	"io/fs"
	"testing"
)

func Test_listShares(t *testing.T) {
	s, err := New("\\\\PHYSREG\\LANELINK", "POS_NCBP", "*Loc!Sms901")
	if err != nil {
		t.Fatal(err)
	}
	err = s.Connect()
	defer s.Close()

	if err != nil {
		t.Fatal(err)
	}
	s.WalkDir("", func(path string, d fs.DirEntry) error {
		fmt.Printf("%s: %s\n", path, d.Name())
		return nil
	})

}

func Test_readFile(t *testing.T) {
	s, err := New("\\\\PHYSREG\\LANELINK", "POS_NCBP", "*Loc!Sms901")
	if err != nil {
		t.Fatal(err)
	}
	err = s.Connect()
	defer s.Close()

	if err != nil {
		t.Fatal(err)
	}
	data, err := s.ReadFile("startup.ini")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
}
