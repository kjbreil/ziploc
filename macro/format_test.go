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
	// test := fmt.Sprintf("testing\nnewline\ncrlf\r\ntext")
	// CRLF([]byte(test))

	fn := "./crlftest.txt"

	// fi, err := os.Stat(fn)
	// if err != nil {
	// 	t.Fatalf("crlftest file does not exist")
	// }

	f, err := os.Open(fn)
	if err != nil {
		t.Fatalf("crlftest file could not be opened")
	}
	//
	// fi, _ := f.Stat()
	//
	// b := make([]byte, fi.Size())
	//
	// _, _ = f.Read(b)
	//
	// b, _ = CRLF(b)
	//
	// o, err := os.Create("./corrected.txt")
	//
	// _, _ = o.Write(b)
	//
	// _, _ = io.Copy(o, f)

	o, err := os.Create("./corrected.txt")

	err = Correct(o, f)
	if err != nil {
		t.Fatal(err)
	}
	t.Fail()

}
