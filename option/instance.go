package option

import (
	"fmt"
	"github.com/kjbreil/ziploc/dss"
	"github.com/kjbreil/ziploc/macro"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func (o *Option) WalkInstance(folder string) error {
	if o.InstanceDir == nil {
		return fmt.Errorf("instanceDir is missing")
	}
	if o.Include == nil {
		return fmt.Errorf("include is missing")
	}
	if o.Exclude == nil {
		return fmt.Errorf("exclude is missing")
	}

	info, err := os.Stat(*o.InstanceDir)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("instance_dir is a file not a folder")
	}

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

func (o *Option) FromInstance() error {
	root := o.root
	for ep, ef := range o.instanceFiles {
		f, err := os.Open(ep)
		if err != nil {
			return err
		}

		var dest string

		m := regexp.MustCompile(`SP-\d{3}`)
		folder := m.FindString(ep)
		if folder == "" {
			folder = o.Folder(ef.Name())
		}

		switch folder {
		case "EVERY":
			dest = o.makeBuildPath(root, dss.Destination(ep), ef.Name())
		case "ROOT":
			dest = o.makeBuildPath(root, "", ef.Name())
		default:
			newPath := filepath.Join(folder, dss.Destination(ep))
			dest = o.makeBuildPath(root, newPath, ef.Name())
		}

		if _, err := os.Stat(filepath.Dir(dest)); os.IsNotExist(err) {
			err = os.MkdirAll(filepath.Dir(dest), 0777)
			if err != nil {
				return fmt.Errorf("could not create directory: %v", err)
			}
		}
		out, err := os.Create(dest)
		if err != nil {
			return err
		}

		err = macro.Correct(out, f, false)
		if err != nil {
			return err
		}

		_ = f.Close()
		_ = out.Close()
	}
	return nil
}
