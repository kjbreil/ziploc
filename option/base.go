package option

import (
	"fmt"
	"github.com/kjbreil/ziploc/dss"
	"github.com/kjbreil/ziploc/macro"
	"log"
	"os"
	"path/filepath"
)

func (o *Option) FromBase(root string) error {
	o.DeleteTemp(root)

	baseWithRoot := filepath.Join(root, o.BaseFolder)

	log.Println("Walking Path", o.BaseFolder)
	// "walk" the BaseFolder for files, adding them to files if the match
	// TODO: REGEX for filename matching (optional?)

	info, err := os.Stat(baseWithRoot)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("base directory not found %v", err)
	}

	err = o.Walk(baseWithRoot)
	if err != nil {
		return err
	}

	// Create a new DSS for this run
	o.dss = dss.New(o.Name, o.Priority)
	// loop over the paths and add to the DSS
	for path, info := range o.files {
		// ignore root files
		if o.Folder(info.Name()) == "ROOT" {
			continue
		}
		err = o.dss.Add(path)
		if err != nil {
			return fmt.Errorf("error adding %s to dss: %v\n", path, err)
		}
	}

	// Write the DSS to the temp directory
	err = o.dss.Write(filepath.Join(root, o.TempDir, o.getType(), o.dss.Name))
	if err != nil {
		return err
	}

	err = o.WriteInstall(root)
	if err != nil {
		return fmt.Errorf("error writing install file: %v\n", err)
	}

	err = o.CopyBase(root)
	if err != nil {
		return err
	}

	err = o.MakeZip(root)
	if err != nil {
		return err
	}

	return nil
}

func (o *Option) CopyBase(root string) error {
	for ep, ef := range o.files {
		f, err := os.Open(ep)
		if err != nil {
			return err
		}
		folder := ep[len(o.BaseFolder)-1 : len(ep)-len(ef.Name())]
		dest := o.makePath(root, folder, ef.Name())

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
