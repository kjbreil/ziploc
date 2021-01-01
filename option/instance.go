package option

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func (o *Option) FromInstance() error {

	if o.InstanceDir == nil {
		return fmt.Errorf("instanceDir is missing")
	}
	if o.Include == nil {
		return fmt.Errorf("include is missing")
	}
	if o.Exclude == nil {
		return fmt.Errorf("exclude is missing")
	}

	log.Println("Walking Path", *o.InstanceDir)

	info, err := os.Stat(*o.InstanceDir)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("instance_dir is a file not a folder")
	}

	err = o.WalkInstance(*o.InstanceDir)
	if err != nil {
		return err
	}

	err = o.CopyFiles("")
	if err != nil {
		return err
	}

	return nil
}

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
			log.Printf("adding file: %s\n", p)
			o.instanceFiles[p] = eachFile
		}
	}
	return nil
}
