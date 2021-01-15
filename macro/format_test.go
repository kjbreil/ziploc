package macro

import (
	"os"
	"testing"
)

// func TestPrimary(t *testing.T) {
// 	type args struct {
// 		function  string
// 		arguments []string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want []byte
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := Primary(tt.args.function, tt.args.arguments...); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Primary() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestPrimary(t *testing.T) {
// 	b, _ := Line(Primary("FMT", "CMP", PrimaryS("WIZGET", "INSTALL_STEP"), "<3", InnserS("FMT", "26")))
// 	// b := InnserS("CHR", "26")
// 	ioutil.WriteFile("./test.txt", b, 0666)
// }

func TestCRLF(t *testing.T) {
	// [35 147 172 63 202 166 188 222 204 205 10]
	// b := []byte("BEFORE #“¬�Ê¦¼ÞÌÍ AFTER")
	//
	// nb, err := CRLF(b)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(string(b))
	// t.Log(string(nb))
	//
	// t.Fail()

	fn := "./crlftest.txt"

	f, err := os.Open(fn)
	if err != nil {
		t.Fatalf("crlftest file could not be opened")
	}

	o, err := os.Create("./corrected.txt")

	err = Correct(o, f, false)
	if err != nil {
		t.Fatal(err)
	}

}
