package dss

import (
	"io"
	"os"

	"github.com/kjbreil/sil/silread"
)

// Read reads DSS data from an io.Reader and returns parsed Entries
func Read(r io.Reader) (Entries, error) {
	var b []byte
	var err error
	var t []Table

	b, err = io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	t = make([]Table, 0)
	err = silread.Unmarshal(b, &t)
	if err != nil {
		return nil, err
	}

	rt := make(Entries, 0, len(t))

	for _, v := range t {
		rt = append(rt, &v)
	}

	return rt, nil
}

// ReadFile reads DSS data from a file and returns parsed Entries
func ReadFile(filename string) (Entries, error) {
	var f *os.File
	var err error

	f, err = os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return Read(f)
}
