package dss

import (
	"github.com/kjbreil/sil/silread"
	"os"
)

func Read(filename string) (Entries, error) {
	t := make([]Table, 0)
	b, _ := os.ReadFile(filename)
	err := silread.Unmarshal(b, &t)
	if err != nil {
		return nil, err
	}

	rt := make(Entries, 0, len(t))

	for _, v := range t {
		rt = append(rt, &v)
	}

	return rt, nil
}
