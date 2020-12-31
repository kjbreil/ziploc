package option

import (
	"fmt"
	"github.com/kjbreil/ziploc/dss"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func (o *Option) WalkInstance(folder string) error {
	if o.excluded(filepath.Base(folder)) {
		return nil
	}

	stat, err := os.Stat(folder)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return fmt.Errorf("%s is not a directory", folder)
	}

	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return err
	}

	for _, eachFile := range files {
		p := filepath.Join(folder, eachFile.Name())
		if eachFile.IsDir() {
			err = o.WalkInstance(p)
			if err != nil {
				return err
			}
		}
		if !eachFile.IsDir() && o.included(eachFile.Name()) && !o.excluded(eachFile.Name()) {
			// log.Printf("adding file: %s\n", p)
			o.Files[p] = eachFile
		}
	}
	return nil
}

func (o *Option) fromInstance() {
	log.Println("Walking Path", o.BaseFolder)
	// "walk" the BaseFolder for files, adding them to files if the match
	// TODO: REGEX for filename matching (optional?)

	info, err := os.Stat(o.BaseFolder)
	if err != nil || !info.IsDir() {
		log.Panicf("Panicing")
	}

	err = o.WalkInstance(o.BaseFolder)

	if err != nil {
		log.Panic(err)
	}

	// Create a new DSS for this run
	o.Dss = dss.New(o.Name, o.Priority)
	// loop over the paths and add to the DSS
	for path, info := range o.Files {
		// ignore root files
		if o.Folder(info.Name()) == "ROOT" {
			continue
		}
		err = o.Dss.Add(path)
		if err != nil {
			log.Panicf("error adding %s to dss: %v\n", path, err)
		}
	}

	// Write the DSS to the temp directory
	err = o.Dss.Write(filepath.Join(o.TempDir, o.getType(), o.Dss.Name))
	if err != nil {
		log.Println(err)
	}

	err = o.WriteInstall()
	if err != nil {
		log.Panicf("error writing install file: %v\n", err)
	}

	err = o.CopyFiles()
	if err != nil {
		log.Panic(err)
	}

	err = o.MakeZip()
	if err != nil {
		log.Panic(err)
	}
}
