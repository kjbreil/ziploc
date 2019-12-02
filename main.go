package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kjbreil/ziploc/dss"
)

type Option struct {
	ZipFilename string
	BaseFolder  string
	Include     []string
	Type        string // Options or Samples
	Files       map[string]os.FileInfo
	DSS         *dss.DSS
}

func main() {

	// setup Option information
	o := Option{
		BaseFolder: "./example/Samples",
		Files:      make(map[string]os.FileInfo),
		Type:       "Samples",
		DSS:        dss.New("PCC SH BASE LB"),
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

	// loop over the paths and
	for path, _ := range o.Files {
		o.DSS.Add(path)
	}

	err = o.DSS.Write(filepath.Join("./out", o.Type, o.DSS.Name))
	if err != nil {
		log.Println(err)
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
