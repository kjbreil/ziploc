package dss

import (
	"github.com/kjbreil/sil/silread"
	"os"
)

func Read(filename string) (Entries, error) {
	var t Entries
	b, _ := os.ReadFile(filename)
	err := silread.Unmarshal(b, &t)
	if err != nil {
		return nil, err
	}

	return t, nil
}
