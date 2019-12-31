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
	Name       string              `json:"name"`
	BaseFolder string              `json:"base_folder"`
	Include    map[string][]string `json:"include"`
	Type       string              `json:"type"` // Options or Samples
	TempDir    string              `json:"temp_dir"`
	ZipOut     string              `json:"zip_out"`
	// Not Exported
	files map[string]os.FileInfo
	dss   *dss.DSS
}

var (
	configLocation = flag.String("c", "", "Config file")
	makeTemplate   = flag.Bool("template", false, "output a template config file and exit")
)

func main() {

	// Parse the flags
	flag.Parse()

	if *makeTemplate {
		configTemplate()
		return
	}

	o, err := readConfig()
	if err != nil {
		log.Panic(err)
	}

	// TODO: Delete the temp directory before doing anything
	// stat the tempdir and delete if there is no error (directory exists)
	_, err = os.Stat(o.TempDir)
	if err == nil {
		log.Printf("Removing temp directory: %s\n", o.TempDir)
		os.RemoveAll(o.TempDir)
	}

	// TODO: Replace spaces in foldername to make safer and easier to work with

	log.Println("Walking Path", o.BaseFolder)
	// "walk" the BaseFolder for files, adding them to files if the match
	// TODO: REGEX for filename matching (optional?)

	info, err := os.Stat(o.BaseFolder)
	if err != nil || !info.IsDir() {
		log.Panicf("Panicing")
	}

	// TODO: change to own code so walk ignores certain folders completely
	err = filepath.Walk(o.BaseFolder, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && o.included(info.Name()) && !excluded(path) {

			log.Printf("adding file: %s\n", path)

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
	for path := range o.files {
		err = o.dss.Add(path)
		if err != nil {
			log.Panicf("error adding %s to dss: %v\n", path, err)
		}
	}

	// Write the DSS to the temp directory
	err = o.dss.Write(filepath.Join(o.TempDir, o.Type, o.dss.Name))
	if err != nil {
		log.Println(err)
	}

	err = o.writeInstall()
	if err != nil {
		log.Panicf("error writing install file: %v\n", err)
	}

	err = o.copyDSSFiles()
	if err != nil {
		log.Panic(err)
	}

	err = o.makeZip()
	if err != nil {
		log.Panic(err)
	}

}

func excluded(current string) bool {
	switch {
	case strings.Contains(strings.ToUpper(filepath.Dir(current)), "TEMP"):
		return true
	case strings.Contains(strings.ToUpper(filepath.Dir(current)), "SAMPLES"):
		return true
	case strings.Contains(strings.ToUpper(filepath.Dir(current)), "OPTIONS"):
		return true
	}

	return false
}

func (o *Option) included(current string) bool {
	b, _ := o.check(current)
	return b
}

func (o *Option) folder(current string) string {
	_, s := o.check(current)
	return s
}

func (o *Option) check(current string) (bool, string) {
	for folder := range o.Include {
		for _, item := range o.Include[folder] {
			// check if it matches, uppercase both
			// fmt.Println(current, item, strings.EqualFold(current, item))
			if strings.EqualFold(current, item) {
				return true, folder
			}
		}
	}
	return false, ""
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
		var dest string
		folder := o.folder(ef.Name())
		switch folder {
		case "every":
			dest = o.makePath(dss.Destination(ep), ef.Name())
		case "root":
			dest = o.makePath("", ef.Name())
		default:
			newPath := filepath.Join(folder, dss.Destination(ep))
			dest = o.makePath(newPath, ef.Name())
		}

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
