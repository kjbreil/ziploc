package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kjbreil/ziploc/dss"
)

type Option struct {
	Name       string   `json:"name,omitempty"`
	BaseFolder string   `json:"base_folder,omitempty"`
	Include    []string `json:"include,omitempty"`
	Type       string   `json:"type,omitempty"` // Options or Samples
	TempDir    string   `json:"temp_dir,omitempty"`
	// Not Exported
	files map[string]os.FileInfo
	dss   *dss.DSS
}

var (
	configLocation = flag.String("c", "", "Config file")
)

func main() {

	// Parse the flags
	flag.Parse()

	o, err := readConfig()
	if err != nil {
		log.Panic(err)
	}

	// TODO: Delete the temp directory before doing anything

	// "walk" the BaseFolder for files, adding them to files if the match
	// TODO: REGEX for filename matching (optional?)
	err = filepath.Walk(o.BaseFolder, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && o.included(info.Name()) {
			o.files[path] = info
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	// Create a new DSS for this run
	o.dss = dss.New(o.Name)
	// loop over the paths and add to the DSS
	for path, _ := range o.files {
		o.dss.Add(path)
	}

	// Write the DSS to the temp directory
	err = o.dss.Write(filepath.Join(o.TempDir, o.Type, o.dss.Name))
	if err != nil {
		log.Println(err)
	}

	o.writeInstall()
	err = o.copyDSSFiles()
	if err != nil {
		log.Panic(err)
	}

	err = o.makeZip()
	if err != nil {
		log.Panic(err)
	}

	configTemplate()

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

// makePath uses filepath.Join to safely create the path to the file using OS independent paths
func (o *Option) makePath(folder string, filename string) string {
	p := filepath.Join(o.TempDir, o.Type, o.dss.Name, folder, filename)
	return p
}

func (o *Option) copyDSSFiles() error {
	for ep, ef := range o.files {
		f, err := os.Open(ep)
		if err != nil {
			return err
		}
		dest := o.makePath(dss.Destination(ep), ef.Name())
		if _, err := os.Stat(filepath.Dir(dest)); os.IsNotExist(err) {
			err = os.MkdirAll(filepath.Dir(dest), 0777)
			if err != nil {
				return fmt.Errorf("could not create directory: %v", err)
			}
		}

		o, err := os.Create(dest)
		if err != nil {
			return err
		}
		_, err = io.Copy(o, f)
		if err != nil {
			return err
		}
		_ = f.Close()
		_ = o.Close()
	}
	return nil
}
