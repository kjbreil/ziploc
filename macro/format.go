package macro

import (
	"golang.org/x/text/encoding/charmap"
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
