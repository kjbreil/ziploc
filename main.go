package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kjbreil/crcloc"
	"github.com/kjbreil/sil"
)

type Option struct {
	ZipFilename string
	BaseFolder  string
	Include     []string
	Files       map[string]os.FileInfo
	DSS         sil.SIL
}

func main() {

	// setup Option information
	o := Option{
		BaseFolder: "./example/Samples",
		Files:      make(map[string]os.FileInfo),
	}
	// just for now until config is built
	o.Include = []string{
		"system.ini",
		"samples.ini",
	}

	// "walk" the BaseFolder for files
	err := filepath.Walk(o.BaseFolder, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && o.included(info.Name()) {
			o.Files[path] = info
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	for path, info := range o.Files {
		b, err := ioutil.ReadFile(path)
		if err != nil {
			log.Panic(err)
		}
		hash := crcloc.Hash(b)
		fmt.Printf("file: %s has hash: %s\n", info.Name(), hash)

	}

}

func (o *Option) included(current string) bool {
	for _, ei := range o.Include {
		// check if it matches, uppercase both
		if strings.ToUpper(current) == strings.ToUpper(ei) {
			return true
		}
	}
	return false

}
