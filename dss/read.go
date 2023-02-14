package dss

import (
	"github.com/kjbreil/sil/silread"
	"log"
	"os"
)

func Read(filename string) *DSS {

	var t Table

	b, _ := os.ReadFile(filename)

	m, err := silread.Unmarshal(b, t)

	if err != nil {
		log.Fatal(err)
	}

	bn, err := m.Marshal(false)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(bn))

	return nil
}
