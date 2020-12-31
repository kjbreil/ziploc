package option

import (
	"fmt"
	"github.com/kjbreil/ziploc/dss"
	"github.com/kjbreil/ziploc/macro"
	"log"
	"os"
	"path/filepath"
)

func (o *Option) FromBase() {
	o.DeleteTemp()

	log.Println("Walking Path", o.BaseFolder)
	// "walk" the BaseFolder for files, adding them to files if the match
	// TODO: REGEX for filename matching (optional?)

	info, err := os.Stat(o.BaseFolder)
	if err != nil || !info.IsDir() {
		log.Panicf("Panicing")
	}

	err = o.Walk(o.BaseFolder)

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
	err = o.Dss.Write(filepath.Join(o.TempDir, o.Type, o.Dss.Name))
	if err != nil {
		log.Println(err)
	}

	err = o.WriteInstall()
	if err != nil {
		log.Panicf("error writing install file: %v\n", err)
	}

	err = o.CopyBase()
	if err != nil {
		log.Panic(err)
	}

	err = o.MakeZip()
	if err != nil {
		log.Panic(err)
	}
}

func (o *Option) CopyBase() error {
	for ep, ef := range o.Files {
		f, err := os.Open(ep)
		if err != nil {
			return err
		}
		folder := ep[len(o.BaseFolder)-1 : len(ep)-len(ef.Name())]
		dest := o.makePath(folder, ef.Name())

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

		err = macro.Correct(out, f)
		if err != nil {
			log.Println(dest)
			return err
		}

		_ = f.Close()
		_ = out.Close()
	}
	return nil
}
