package macro

import (
	"golang.org/x/text/encoding/charmap"
	"os"
)

// Line returns a line encoded in Windows1252 with a line ending
func Line(b []byte) ([]byte, error) {
	b = append(b, []byte("\r\n")...)
	return charmap.Windows1252.NewEncoder().Bytes(b)
}

// LineS returns a line encoded in Windows1252 with a line ending in a string
// IGNORES ANY ERROR
func LineS(s string) string {
	b, _ := Line([]byte(s))
	return string(b)
}

// CRLF returns a new byte slice that has been converted
func CRLF(b []byte) ([]byte, error) {
	for i := range b {
		// check if first character is newline, convert if needed
		if i == 0 && b[i] == 10 {
			b = append(b, 0)
			copy(b[i+1:], b[i:])
			b[i] = 13
			continue
		}
		if b[i] == 10 && b[i-1] != 13 {
			// fmt.Println("newline found")
			b = append(b, 0)
			copy(b[i+1:], b[i:])
			b[i] = 13
		}

	}
	return charmap.Windows1252.NewEncoder().Bytes(b)
}

func Correct(dst *os.File, src *os.File) error {
	info, err := src.Stat()
	if err != nil {
		return err
	}
	b := make([]byte, info.Size())
	_, err = src.Read(b)
	if err != nil {
		return err
	}
	b, err = CRLF(b)
	if err != nil {
		return err
	}
	_, err = dst.Write(b)
	if err != nil {
		return err
	}
	return nil
}
