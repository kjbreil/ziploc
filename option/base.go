package option

import (
	"fmt"
	"github.com/kjbreil/ziploc/dss"
	"github.com/kjbreil/ziploc/iniUpdater"
	"github.com/kjbreil/ziploc/macro"
	"log"
	"os"
	"path/filepath"
)

func (o *Option) GetBaseFiles() error {
	baseWithRoot := filepath.Join(o.root, o.BaseFolder)

	info, err := os.Stat(baseWithRoot)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("base directory not found %v", err)
	}
	err = o.Walk(baseWithRoot)
	if err != nil {
		return err
	}
	return nil
}

func (o *Option) GetBaseDSS() error {
	// Create a new DSS for this run
	o.dss = dss.New(o.Name, o.Priority)
	// loop over the paths and add to the DSS
	for path, info := range o.files {
		switch {
		// ignore root files
		case o.Folder(info.Name()) == "ROOT":
			continue
		// ignore unix hidden files
		case info.Name()[0:1] == ".":
			continue
		}
		err := o.dss.Add(path)
		if err != nil {
			return fmt.Errorf("error adding %s to dss: %v\n", path, err)
		}
	}

	return nil
}

func (o *Option) FromBase(root string) error {
	o.DeleteTemp(root)
	var err error

	log.Println("Walking Path", o.BaseFolder)
	// "walk" the BaseFolder for files, adding them to files if the match
	// TODO: REGEX for filename matching (optional?)

	// o.GetBaseFiles()

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
	log.Println("-->", root)
	// b := filepath.Join(root, o.BaseFolder)
	for ep, ef := range o.files {
		var force bool // to force correction, only when its an ini for now
		f, err := os.Open(ep)
		if err != nil {
			return err
		}
		// folder := ep[len(b) : len(ep)-len(ef.Name())]
		folder := filepath.Base(filepath.Dir(ep))
		dest := o.makePath(root, folder, ef.Name())

		// is the file a part of the ini map, maybe do something then
		noExt := ef.Name()[:len(ef.Name())-len(filepath.Ext(ef.Name()))]
		ini, ok := o.IniMaps[noExt]
		if ok {
			if ini.Original != nil {
				f, err = iniUpdater.File(*ini.Original, *ini.Base, folder)
				if err != nil {
					return err
				}
				dest = o.makePath(root, "INBOX", noExt+"_INI.SQL")
				force = true
			}
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

		err = macro.Correct(out, f, force)
		if err != nil {
			return err
		}

		_ = f.Close()
		_ = out.Close()
	}
	return nil
}
