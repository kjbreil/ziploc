package macro

import (
	"fmt"
	"os"
	"path"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
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
func CRLF(ob []byte) ([]byte, error) {

	for i := range ob {
		// check if first character is newline, convert if needed
		if i == 0 && ob[i] == 10 {
			ob = append(ob, 0)
			copy(ob[i+1:], ob[i:])
			ob[i] = 13
			continue
		}
		if ob[i] == 10 && ob[i-1] != 13 {
			ob = append(ob, 0)
			copy(ob[i+1:], ob[i:])
			ob[i] = 13
		}
	}

	enc := charmap.Windows1252.NewEncoder()
	var i int
	var fb []byte
	// dis is magic
	// basically loop over ob using the enc.Transform to put bytes into a final object until you cannot encode something
	// then just pop that cannot be encoded byte to the end but that doesn't sound right
	// shouldn't i just discard the problem byte and call it good, ala registered symbol having a prepended byte in
	// utf-8 but this works so fuck it
	for i < len(ob) {
		ib := make([]byte, len(ob)-i)
		di, si, err := enc.Transform(ib, ob[i:], false)
		// we hit an error, take where we are
		if err != nil {
			fb = append(fb, ib[:di]...)
			fb = append(fb, ob[di+i])
			if err == transform.ErrShortSrc {
				break
			}
		} else {
			fb = append(fb, ib[:di]...)

		}

		i = si + i + 1
	}
	return fb, nil
}

func Correct(dst *os.File, src *os.File, force bool) error {
	extToCorrect := map[string]bool{
		".htm": true,
		".sqi": true,
		".sql": true,
		".ini": true,
		".txt": true,
	}

	info, err := src.Stat()
	if err != nil {
		return err
	}
	b := make([]byte, info.Size())
	_, err = src.Read(b)
	if err != nil {
		return fmt.Errorf("read file error %v", err)
	}

	_, ok := extToCorrect[strings.ToLower(path.Ext(src.Name()))]
	if ok || force {
		b, err = CRLF(b)
		if err != nil {
			return fmt.Errorf("CRLF conversion failed: %v", err)
		}
	}

	_, err = dst.Write(b)
	if err != nil {
		return err
	}
	return nil
}
